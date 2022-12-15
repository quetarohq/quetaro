package quetaro

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/quetarohq/quetaro/ctxutil"
	"github.com/rs/zerolog/log"
)

type IntakeInvokeAgent struct {
	*IntakeInvoke
}

func newIntakeInvokeAgent(intakeInvoke *IntakeInvoke) *IntakeInvokeAgent {
	agent := &IntakeInvokeAgent{
		IntakeInvoke: intakeInvoke,
	}

	return agent
}

func (agent *IntakeInvokeAgent) run(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Info().Msg("start agent")

	err := loopForAgent(ctx, agent.ConnConfig, agent.Interval, agent.ErrInterval, func(ctx context.Context, conn *pgx.Conn) error {
		err := agent.invoke(ctx, conn)

		if err != nil {
			return errors.Wrap(err, "failed to invoke a job")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "cannot continue to run agent")
	}

	return nil
}

func (agent *IntakeInvokeAgent) invoke(ctx context.Context, conn *pgx.Conn) error {
	// do not cancel the following processes
	ctx = ctxutil.WithoutCancel(ctx)

	err := pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		var msgId, funcName, payload string
		sql, args := sq.Select("id", "function_name", "payload").From("jobs").
			Where(sq.And{
				sq.Eq{"queue_name": agent.QueueName},
				sq.Or{
					sq.Eq{"status": JobStatusPending},
					sq.Eq{"status": JobStatusFailure},
				},
				sq.LtOrEq{
					"invoke_after": time.Now(),
				},
			},
			).OrderBy("invoke_after", "updated_at").Limit(1).
			Suffix("for update skip locked").MustSql()
		err := tx.QueryRow(ctx, sql, args...).Scan(&msgId, &funcName, &payload)

		if err == pgx.ErrNoRows {
			return nil
		} else if err != nil {
			return errors.Wrap(err, "failed to fetch a job from DB")
		}

		logger := log.Ctx(ctx).With().Str("id", msgId).Str("function_name", funcName).Logger()
		invokeErr := awsutil.InvokeFunction(ctx, agent.lambda, funcName, payload, map[string]string{JobIdKey: msgId})

		valByCol := sq.Eq{
			"status":     JobStatusInvoked,
			"updated_at": time.Now(),
		}

		if invokeErr != nil {
			valByCol["status"] = JobStatusInvokeFailure
			valByCol["last_error"] = invokeErr.Error()
		}

		sql, args = sq.Update("jobs").SetMap(valByCol).Where(sq.Eq{"id": msgId}).MustSql()
		_, err = tx.Exec(ctx, sql, args...)

		if err != nil {
			ewj := &ErrWithJob{cause: err, Id: msgId, Name: funcName}
			return errors.Wrap(ewj, "failed to update the job status")
		}

		if invokeErr == nil {
			logger.Info().Msg("lambda function was invoked")
		} else {
			logger.Error().Err(invokeErr).Msg("failed to invoke the lambda function")
		}

		return nil
	}) // end of pgx.BeginFunc()

	if err != nil {
		return errors.Wrap(err, "error in agant transactions")
	}

	return nil
}
