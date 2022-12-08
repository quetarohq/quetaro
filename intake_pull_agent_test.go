package quetaro

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_IntakePullAgent_run(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	agent := &IntakePullAgent{
		IntakePull: &IntakePull{
			IntakePullOpts: &IntakePullOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-intake-test",
		},
		sqs: res.SQS,
	}

	testutil.SendMessage(t, res.SQS, "qtr-intake-test", "qtr-job-test", `{"foo":"bar","zoo":true}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs")
	assert.Len(rows, 1)
	r := rows[0]

	for k, v := range map[string]any{
		"queue_name":    "qtr-intake-test",
		"function_name": "qtr-job-test",
		"payload":       map[string]any{"foo": "bar", "zoo": true},
		"status":        "pending",
	} {
		assert.Equal(v, r[k])
	}

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-intake-test")
	assert.Len(msgs, 0)
}

func Test_IntakePullAgent_run_NoFuncName(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	agent := &IntakePullAgent{
		IntakePull: &IntakePull{
			IntakePullOpts: &IntakePullOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-intake-test",
		},
		sqs: res.SQS,
	}

	testutil.SendMessage(t, res.SQS, "qtr-intake-test", "", `{"foo":"bar","zoo":true}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs")
	assert.Len(rows, 1)
	r := rows[0]

	for k, v := range map[string]any{
		"queue_name":    "qtr-intake-test",
		"function_name": "",
		"payload":       map[string]any{"foo": "bar", "zoo": true},
		"status":        "invalid",
	} {
		assert.Equal(v, r[k])
	}

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-intake-test")
	assert.Len(msgs, 0)
}

func Test_IntakePullAgent_run_InvalidJSON(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	agent := &IntakePullAgent{
		IntakePull: &IntakePull{
			IntakePullOpts: &IntakePullOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-intake-test",
		},
		sqs: res.SQS,
	}

	testutil.SendMessage(t, res.SQS, "qtr-intake-test", "qtr-job-test", `xxx`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)

	rows := testutil.Query(t, res.Conn, "select * from jobs")
	assert.Len(rows, 1)
	r := rows[0]

	for k, v := range map[string]any{
		"queue_name":    "qtr-intake-test",
		"function_name": "qtr-job-test",
		"payload":       map[string]any{},
		"status":        "invalid",
	} {
		assert.Equal(v, r[k])
	}

	msgs := testutil.ReceiveMessages(t, res.SQS, "qtr-intake-test")
	assert.Len(msgs, 0)
}

func Test_IntakePullAgent_run_DeleteMessages_failed(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	var called bool
	agent := &IntakePullAgent{
		IntakePull: &IntakePull{
			IntakePullOpts: &IntakePullOpts{
				QueueName:  "qtr-intake-test",
				ConnConfig: res.ConnCfg,
				MaxRecvNum: 10,
			},
			QueueUrl: "qtr-intake-test",
		},
		sqs: &testutil.MockSQS{
			SQS: res.SQS,
			MockDeleteMessageBatch: func(ctx context.Context, dmbi *sqs.DeleteMessageBatchInput, f ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error) {
				called = true
				return &sqs.DeleteMessageBatchOutput{
					Failed: []types.BatchResultErrorEntry{
						{
							Code:    aws.String("Code"),
							Id:      dmbi.Entries[0].Id,
							Message: aws.String("DeleteMessageBatch error"),
						},
					},
				}, errors.New("DeleteMessageBatch error")
			},
		},
	}

	testutil.SendMessage(t, res.SQS, "qtr-intake-test", "qtr-job-test", `{"foo":"bar","zoo":true}`)

	err := agent.run(&testutil.MockCtx{F: "loop_for_agent.go", N: 2})
	assert.NoError(err)
	assert.True(called)

	rows := testutil.Query(t, res.Conn, "select * from jobs")
	assert.Equal(0, len(rows))
}
