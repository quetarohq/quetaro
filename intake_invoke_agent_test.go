package quetaro

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
	"github.com/winebarrel/quetaro/internal/testutil"
)

func Test_IntakeInvokeAgent_run_Success(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

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
			'pending',
			now(),
			0,
			now(),
			now()
		)
`)

	agent := &IntakeInvokeAgent{
		IntakeInvoke: &IntakeInvoke{
			IntakeInvokeOpts: &IntakeInvokeOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
			},
			lambda: testutil.NewLambdClient(t),
		},
	}

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Equal(1, len(rows))
	r := rows[0]
	assert.Equal("invoked", r["status"])

	successMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-success-test")
	assert.Len(successMsgs, 1)
	m := successMsgs[0]
	j, _ := fastjson.Parse(aws.ToString(m.Body))
	assert.Equal("013eb466-184c-43e6-b0c2-6667d5cf3b47", string(j.GetStringBytes("requestPayload", "_id")))

	failureMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(failureMsgs, 0)
}

func Test_IntakeInvokeAgent_run_Failure(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

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
	'{"_fail":true}',
	'pending',
	now(),
	0,
	now(),
	now()
)
`)

	agent := &IntakeInvokeAgent{
		IntakeInvoke: &IntakeInvoke{
			IntakeInvokeOpts: &IntakeInvokeOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
			},
			lambda: testutil.NewLambdClient(t),
		},
	}

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Equal(1, len(rows))
	r := rows[0]
	assert.Equal("invoked", r["status"])

	successMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-success-test")
	assert.Len(successMsgs, 0)

	failureMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(failureMsgs, 1)
	m := failureMsgs[0]
	j, _ := fastjson.Parse(aws.ToString(m.Body))
	assert.Equal("013eb466-184c-43e6-b0c2-6667d5cf3b47", string(j.GetStringBytes("requestPayload", "_id")))
}

func Test_IntakeInvokeAgent_run_InvokeError(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

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
	'{"_fail":true}',
	'pending',
	now(),
	0,
	now(),
	now()
)
`)

	var called bool
	agent := &IntakeInvokeAgent{
		IntakeInvoke: &IntakeInvoke{
			IntakeInvokeOpts: &IntakeInvokeOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
			},
			lambda: &testutil.MockLambda{
				Lambda: testutil.NewLambdClient(t),
				MockInvoke: func(ctx context.Context, ii *lambda.InvokeInput, f ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
					called = true
					return nil, errors.New("Invoke error")
				},
			},
		},
	}

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)
	assert.True(called)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Equal(1, len(rows))
	r := rows[0]
	assert.Equal("invoke_failure", r["status"])
	assert.Equal("Lambda Invoke error: Invoke error", r["last_error"])

	successMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-success-test")
	assert.Len(successMsgs, 0)

	failureMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(failureMsgs, 0)
}

func Test_IntakeInvokeAgent_run_Skip(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

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
	'{"_fail":true}',
	'pass',
	now(),
	0,
	now(),
	now()
)
`)

	agent := &IntakeInvokeAgent{
		IntakeInvoke: &IntakeInvoke{
			IntakeInvokeOpts: &IntakeInvokeOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
			},
			lambda: testutil.NewLambdClient(t),
		},
	}

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Equal(1, len(rows))
	r := rows[0]
	assert.Equal("pass", r["status"])

	successMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-success-test")
	assert.Len(successMsgs, 0)

	failureMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(failureMsgs, 0)
}

func Test_IntakeInvokeAgent_run_NotRunYet(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

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
	'{"_fail":true}',
	'pending',
	now() + cast('1 year' as interval),
	0,
	now(),
	now()
)
`)

	agent := &IntakeInvokeAgent{
		IntakeInvoke: &IntakeInvoke{
			IntakeInvokeOpts: &IntakeInvokeOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
			},
			lambda: testutil.NewLambdClient(t),
		},
	}

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs where id = '013eb466-184c-43e6-b0c2-6667d5cf3b47'")
	assert.Equal(1, len(rows))
	r := rows[0]
	assert.Equal("pending", r["status"])

	successMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-success-test")
	assert.Len(successMsgs, 0)

	failureMsgs := testutil.ReceiveMessages(t, res.SQS, "qtr-outlet-failure-test")
	assert.Len(failureMsgs, 0)
}
