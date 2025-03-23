package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var (
	DatabaseDsn = "postgres://qtr_test@127.0.0.1:15432?sslmode=disable"
)

func init() {
	if dsn := os.Getenv("TEST_QTR_DATABASE_DSN"); dsn != "" {
		DatabaseDsn = dsn
	}
}

func NewConnConfig(t *testing.T) *pgx.ConnConfig {
	t.Helper()
	connCfg, err := pgx.ParseConfig(DatabaseDsn)

	if err != nil {
		panic(err)
	}

	return connCfg
}

func ConnectDB(t *testing.T, connCfg *pgx.ConnConfig) *pgx.Conn {
	t.Helper()
	conn, err := pgx.ConnectConfig(context.Background(), connCfg)

	if err != nil {
		panic(err)
	}

	return conn
}

func Query(t *testing.T, conn *pgx.Conn, sql string, args ...any) []map[string]any {
	t.Helper()
	rows, err := conn.Query(context.Background(), sql, args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	fd := rows.FieldDescriptions()
	mapRows := []map[string]any{}

	for rows.Next() {
		m := map[string]any{}
		vals, err := rows.Values()

		if err != nil {
			panic(err)
		}

		for i, v := range vals {
			col := fd[i].Name
			m[col] = v
		}

		mapRows = append(mapRows, m)
	}

	return mapRows
}
