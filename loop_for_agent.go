package quetaro

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro/dbutil"
	"github.com/rs/zerolog/log"
)

func loopForAgent(ctx context.Context, connCfg *pgx.ConnConfig, interval time.Duration, errInterval time.Duration, proc func(context.Context, *pgx.Conn) error) error {
	logger := log.Ctx(ctx)
	conn, err := dbutil.Connect(ctx, connCfg)

	if err != nil {
		return errors.Wrap(err, "failed to connect DB")
	}

	defer func() {
		if !conn.IsClosed() {
			conn.Close(ctx)
		}
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
			// do nothing
		}

		err := proc(ctx, conn)

		if err != nil {
			logger.Error().Stack().Err(err).Msg("error in agent loop")

			select {
			case <-time.After(errInterval):
				// do nothing
			case <-ctx.Done():
				break LOOP
			}

			if !conn.IsClosed() {
				conn.Close(ctx)
			}

			conn, err = dbutil.Connect(ctx, connCfg)

			if err != nil {
				return errors.Wrap(err, "failed to reconnect DB")
			}
		} else {
			select {
			case <-time.After(interval):
				// do nothing
			case <-ctx.Done():
				break LOOP
			}
		}
	}

	return nil
}
