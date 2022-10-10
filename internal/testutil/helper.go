package testutil

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackc/pgx/v5"
)

type Resource struct {
	Conn    *pgx.Conn
	ConnCfg *pgx.ConnConfig
	SQS     *sqs.Client
}

func (r *Resource) Cleanup(t *testing.T) {
	t.Helper()
	r.Conn.Close(context.Background())
	PurgeQueue(t, r.SQS, "qtr-intake-test")
	PurgeQueue(t, r.SQS, "qtr-outlet-success-test")
	PurgeQueue(t, r.SQS, "qtr-outlet-failure-test")
}

func SetupAgent(t *testing.T) *Resource {
	t.Helper()
	connCfg := NewConnConfig(t)
	conn := ConnectDB(t, connCfg)
	_, err := conn.Exec(context.Background(), "delete from jobs")

	if err != nil {
		panic(err)
	}

	return &Resource{
		Conn:    conn,
		ConnCfg: connCfg,
		SQS:     NewSQSClient(t),
	}
}
