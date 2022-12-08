package quetaro

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/quetarohq/quetaro/ctxutil"
	"github.com/rs/zerolog/log"
)

const (
	FuncAttreName = "FunctionName"
)

type IntakePullAgent struct {
	*IntakePull
	sqs awsutil.SQSAPI
}

func newIntakePullAgent(intakePull *IntakePull) *IntakePullAgent {
	sqsClient := sqs.NewFromConfig(intakePull.AwsCfg)

	agent := &IntakePullAgent{
		IntakePull: intakePull,
		sqs:        sqsClient,
	}

	return agent
}

func (agent *IntakePullAgent) run(ctx context.Context) error {
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

func (agent *IntakePullAgent) pull(ctx context.Context, conn *pgx.Conn) error {
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
			msgId := aws.ToString(m.MessageId)
			logger := logger.With().Str("id", msgId).Logger()

			var alreadyExists bool
			sql, args := sq.Select("*").Prefix("select exists (").From("jobs").
				Where(sq.Eq{"id": msgId}).Limit(1).Suffix(")").MustSql()
			err := tx.QueryRow(ctx, sql, args...).Scan(&alreadyExists)

			if err != nil {
				ewj := &ErrWithJob{cause: err, Id: msgId}
				return errors.Wrap(ewj, "failed to fetch a job from DB")
			}

			if alreadyExists {
				logger.Warn().Msg("message already queued")
				msgsToDel = append(msgsToDel, m)
				continue
			}

			var funcName string
			var invalidCause []string

			if funcAttr, ok := m.MessageAttributes[FuncAttreName]; ok {
				funcName = aws.ToString(funcAttr.StringValue)
			} else {
				invalidCause = append(invalidCause, "no function name")
			}

			var body string
			rawBody := aws.ToString(m.Body)

			if json.Valid([]byte(rawBody)) {
				body = rawBody
				logger = logger.With().Int("body_size", len(body)).Logger()
			} else {
				invalidCause = append(invalidCause, "invalid JSON")
			}

			isValid := funcName != "" && body != ""
			now := time.Now()

			valByCol := sq.Eq{
				"id":            msgId,
				"queue_name":    agent.QueueName,
				"function_name": funcName,
				"payload":       body,
				"status":        JobStatusPending,
				"invoke_after":  now,
				"error_count":   0,
				"created_at":    now,
				"updated_at":    now,
			}

			if body == "" {
				valByCol["payload"] = "{}"
			}

			if !isValid {
				valByCol["status"] = JobStatusInvalid
				valByCol["last_error"] = strings.Join(invalidCause, ",")
			}

			sql, args = sq.Insert("jobs").SetMap(valByCol).MustSql()
			_, err = tx.Exec(ctx, sql, args...)

			if err != nil {
				ewj := &ErrWithJob{cause: err, Id: msgId, Name: funcName}
				return errors.Wrap(ewj, "failed to insert the message to DB")
			}

			if isValid {
				logger = logger.With().Str("function_name", funcName).Logger()
				logger.Info().Msg("message queued")
			} else {
				logger.Error().Strs("invalid_cause", invalidCause).Msg("invalid message received")
			}

			msgsToDel = append(msgsToDel, m)
		}

		if len(msgsToDel) > 0 {
			err, failed := awsutil.DeleteMessages(ctx, agent.sqs, agent.QueueUrl, msgsToDel)

			if err != nil {
				for _, e := range failed {
					logger.Error().EmbedObject(FromBatchResultErrorEntry(e)).Msg("failed to delete the message from SQS")
					errMsgId := aws.ToString(e.Id)
					sql, args := sq.Delete("jobs").Where(sq.Eq{"id": errMsgId}).MustSql()
					_, err = tx.Exec(ctx, sql, args...)

					if err != nil {
						ewj := &ErrWithJob{cause: err, Id: errMsgId}
						return errors.Wrap(ewj, "failed to delete the job on DeleteMessageBatch failure")
					}
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
