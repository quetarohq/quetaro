package testutil

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type MockCtx struct {
	F     string
	N     int
	count int
}

func (*MockCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (m *MockCtx) Done() <-chan struct{} {
	_, file, _, _ := runtime.Caller(1)

	if strings.HasSuffix(file, m.F) {
		m.count++

		if m.count >= m.N {
			done := make(chan struct{}, 1)
			done <- struct{}{}
			return done
		}
	}

	return nil
}

func (*MockCtx) Err() error {
	return nil
}

func (*MockCtx) Value(key any) any {
	return nil
}

type MockSQS struct {
	SQS                    *sqs.Client
	MockReceiveMessage     func(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	MockDeleteMessageBatch func(context.Context, *sqs.DeleteMessageBatchInput, ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error)
}

func (m *MockSQS) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	if m.MockReceiveMessage != nil {
		return m.MockReceiveMessage(ctx, params, optFns...)
	} else {
		return m.SQS.ReceiveMessage(ctx, params, optFns...)
	}
}

func (m *MockSQS) DeleteMessageBatch(ctx context.Context, params *sqs.DeleteMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error) {
	if m.MockDeleteMessageBatch != nil {
		return m.MockDeleteMessageBatch(ctx, params, optFns...)
	} else {
		return m.SQS.DeleteMessageBatch(ctx, params, optFns...)
	}
}

type MockLambda struct {
	Lambda     *lambda.Client
	MockInvoke func(context.Context, *lambda.InvokeInput, ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

func (m *MockLambda) Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	if m.MockInvoke != nil {
		return m.MockInvoke(ctx, params, optFns...)
	} else {
		return m.Lambda.Invoke(ctx, params, optFns...)
	}
}
