package quetaro

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_OutletFailureAgent_run(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	agent := &OutletFailureAgent{
		OutletFailure: &OutletFailure{
			OutletFailureOpts: &OutletFailureOpts{
				QueueName:  "qtr-outlet-failure-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-outlet-failure-test",
		},
		sqs: res.SQS,
	}

	//nolint:errcheck
	res.Conn.Exec(context.Background(), `
		insert into jobs (
			id,
			queue_name,
			function_name,
			payload,
			status,
			invoke_after,
			error_count,
			created_at,
			updated_at
		) values (
			'013eb466-184c-43e6-b0c2-6667d5cf3b47',
			'qtr-intake-test',
			'qtr-job-test',
			'{}',
			'invoked',
			now(),
			0,
			now(),
			now()
		)
		`)

	testutil.SendMessage(t, res.SQS, "qtr-outlet-failure-test", "", `{"requestPayload":{"_id": "013eb466-184c-43e6-b0c2-6667d5cf3b47"}}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Len(rows, 1)
	r := rows[0]
	assert.Equal("failure", r["status"])

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(msgs, 0)
}

func Test_OutletFailureAgent_run_NoRecord(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	agent := &OutletFailureAgent{
		OutletFailure: &OutletFailure{
			OutletFailureOpts: &OutletFailureOpts{
				QueueName:  "qtr-outlet-failure-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-outlet-failure-test",
		},
		sqs: res.SQS,
	}

	testutil.SendMessage(t, res.SQS, "qtr-outlet-failure-test", "", `{"requestPayload":{"_id": "013eb466-184c-43e6-b0c2-6667d5cf3b47"}}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(msgs, 0)
}

func Test_OutletFailureAgent_run_DeleteMessages_failed(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	var called bool
	agent := &OutletFailureAgent{
		OutletFailure: &OutletFailure{
			OutletFailureOpts: &OutletFailureOpts{
				QueueName:  "qtr-outlet-failure-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-outlet-failure-test",
		},
		sqs: &testutil.MockSQS{
			SQS: res.SQS,
			MockDeleteMessageBatch: func(ctx context.Context, dmbi *sqs.DeleteMessageBatchInput, f ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error) {
				called = true
				return &sqs.DeleteMessageBatchOutput{}, errors.New("DeleteMessageBatch error")
			},
		},
	}

	//nolint:errcheck
	res.Conn.Exec(context.Background(), `
		insert into jobs (
			id,
			queue_name,
			function_name,
			payload,
			status,
			invoke_after,
			error_count,
			created_at,
			updated_at
		) values (
			'013eb466-184c-43e6-b0c2-6667d5cf3b47',
			'qtr-intake-test',
			'qtr-job-test',
			'{}',
			'invoked',
			now(),
			0,
			now(),
			now()
		)
		`)

	testutil.SendMessage(t, res.SQS, "qtr-outlet-failure-test", "", `{"requestPayload":{"_id": "013eb466-184c-43e6-b0c2-6667d5cf3b47"}}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)
	assert.True(called)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Len(rows, 1)
	r := rows[0]
	assert.Equal("failure", r["status"])

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(msgs, 0)
}
