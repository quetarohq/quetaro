package dbutil

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func Connect(ctx context.Context, connCfg *pgx.ConnConfig) (*pgx.Conn, error) {
	conn, err := pgx.ConnectConfig(ctx, connCfg)

	if err != nil {
		return nil, errors.Wrap(err, "pgx.ConnectConfig() error")
	}

	err = conn.Ping(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "conn.Ping() error")
	}

	return conn, nil
}
