package quetaro

import (
	"context"

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

type OutletSuccessAgent struct {
	*OutletSuccess
	sqs awsutil.SQSAPI
}

func newOutletSuccessAgent(outletSuccess *OutletSuccess) *OutletSuccessAgent {
	sqsClient := sqs.NewFromConfig(outletSuccess.AwsCfg, func(o *sqs.Options) {
		if outletSuccess.AWSEndpointUrl != "" {
			o.BaseEndpoint = aws.String(outletSuccess.AWSEndpointUrl)
		}
	})

	agent := &OutletSuccessAgent{
		OutletSuccess: outletSuccess,
		sqs:           sqsClient,
	}

	return agent
}

func (agent *OutletSuccessAgent) run(ctx context.Context) error {
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

func (agent *OutletSuccessAgent) pull(ctx context.Context, conn *pgx.Conn) error {
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
			sql, args := sq.Select("id", "function_name").From("jobs").
				Where(sq.Eq{"id": msgJobId}).Limit(1).Suffix("for update").MustSql()
			err = tx.QueryRow(ctx, sql, args...).Scan(&jobId, &funcName)

			if err == pgx.ErrNoRows {
				logger.Error().Msg("no job in queue")
				msgsToDel = append(msgsToDel, m)
				continue
			} else if err != nil {
				ewj := &ErrWithJob{cause: err, Id: msgJobId}
				return errors.Wrap(ewj, "failed to fetch a job from DB")
			}

			logger = logger.With().Str("function_name", funcName).Logger()
			sql, args = sq.Delete("jobs").Where(sq.Eq{"id": jobId}).MustSql()
			_, err = tx.Exec(ctx, sql, args...)

			if err != nil {
				ewj := &ErrWithJob{cause: err, Id: jobId, Name: funcName}
				return errors.Wrap(ewj, "failed to delete the job")
			}

			msgsToDel = append(msgsToDel, m)
			logger.Info().Msg("delete the successful job")
		}

		if len(msgsToDel) > 0 {
			failed, err := awsutil.DeleteMessages(ctx, agent.sqs, agent.QueueUrl, msgsToDel)

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
