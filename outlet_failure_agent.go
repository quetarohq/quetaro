package quetaro

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/quetarohq/quetaro/ctxutil"
	"github.com/quetarohq/quetaro/sqsmsg"
	"github.com/rs/zerolog/log"
)

type OutletFailureAgent struct {
	*OutletFailure
	sqs awsutil.SQSAPI
}

func newOutletFailureAgent(outletFailure *OutletFailure) *OutletFailureAgent {
	sqsClient := sqs.NewFromConfig(outletFailure.AwsCfg)

	agent := &OutletFailureAgent{
		OutletFailure: outletFailure,
		sqs:           sqsClient,
	}

	return agent
}

func (agent *OutletFailureAgent) run(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Info().Msg("start agent")

	err := loopForAgent(ctx, agent.ConnConfig, agent.Interval, agent.ErrInterval, func(ctx context.Context, conn *pgx.Conn) error {
		err := agent.pull(ctx, conn)

		if err != nil {
			return errors.Wrap(err, "failed to pull messages")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "cannot continue to run agent")
	}

	return nil
}

func (agent *OutletFailureAgent) pull(ctx context.Context, conn *pgx.Conn) error {
	msgs, err := awsutil.ReceiveMessages(ctx, agent.sqs, agent.QueueUrl, agent.MaxRecvNum, FuncAttreName)

	if err != nil {
		return errors.Wrap(err, "failed to receive messages")
	}

	if len(msgs) == 0 {
		return nil
	}

	// do not cancel the following processes
	ctx = ctxutil.WithoutCancel(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		logger := log.Ctx(ctx)
		msgsToDel := []types.Message{}

		for _, m := range msgs {
			sqsMsgId := aws.ToString(m.MessageId)
			sqsBody := aws.ToString(m.Body)
			msgJobId, err := sqsmsg.DecodeId(sqsBody, JobIdKey)

			if err != nil {
				logger.Error().Err(err).Str("sqs_message_id", sqsMsgId).Str("body", sqsBody).
					Msg("failed to decode job id")
				msgsToDel = append(msgsToDel, m)
				continue
			}

			logger := logger.With().Str("id", msgJobId).Logger()

			var jobId, funcName string
			var invokeAfter time.Time
			var errCount int
			sql, args := sq.Select("id", "function_name", "invoke_after", "error_count").From("jobs").
				Where(sq.Eq{"id": msgJobId}).Limit(1).Suffix("for update").MustSql()
			err = tx.QueryRow(ctx, sql, args...).Scan(&jobId, &funcName, &invokeAfter, &errCount)

			if err == pgx.ErrNoRows {
				logger.Error().Msg("no job in queue")
				msgsToDel = append(msgsToDel, m)
				continue
			} else if err != nil {
				ewj := &ErrWithJob{cause: err, Id: msgJobId}
				return errors.Wrap(ewj, "failed to fetch a job from DB")
			}

			errCount += 1
			interval := errCount*errCount*errCount*errCount + 3 // errCount^4+３
			logger = logger.With().Str("function_name", funcName).Int("error_count", errCount).Logger()

			valByCol := sq.Eq{
				"status":       JobStatusFailure,
				"error_count":  errCount,
				"invoke_after": invokeAfter.Add(time.Duration(interval) * time.Second),
				"last_error":   sqsBody,
				"updated_at":   time.Now(),
			}

			sql, args = sq.Update("jobs").SetMap(valByCol).Where(sq.Eq{"id": jobId}).MustSql()
			_, err = tx.Exec(ctx, sql, args...)

			if err != nil {
				ewj := &ErrWithJob{cause: err, Id: jobId, Name: funcName}
				return errors.Wrap(ewj, "failed to update the job")
			}

			msgsToDel = append(msgsToDel, m)
			logger.Info().Msg("update the failed job status")
		}

		if len(msgsToDel) > 0 {
			err, failed := awsutil.DeleteMessages(ctx, agent.sqs, agent.QueueUrl, msgsToDel)

			if err != nil {
				for _, e := range failed {
					logger.Error().EmbedObject(FromBatchResultErrorEntry(e)).Msg("failed to delete the message from SQS")
				}

				// continue if DeleteMessageBatch fails
				logger.Error().Err(err).Msgf("failed to delete messages from SQS")
			}
		}

		return nil
	}) // end of pgx.BeginFunc()

	if err != nil {
		return errors.Wrap(err, "error in agant transactions")
	}

	return nil
}
