package awsutil_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInvokeFunction(t *testing.T) {
	assert := assert.New(t)
	var called bool

	mock := &testutil.MockLambda{
		MockInvoke: func(ctx context.Context, ii *lambda.InvokeInput, f ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
			called = true
			assert.Equal("funcName", aws.ToString(ii.FunctionName))
			assert.Equal(types.InvocationTypeEvent, ii.InvocationType)
			assert.Equal([]byte(`{"_id":"ida","foo":"bar"}`), ii.Payload)
			return &lambda.InvokeOutput{}, nil
		},
	}

	err := awsutil.InvokeFunction(context.Background(), mock, "funcName", `{"foo":"bar"}`, map[string]string{"_id": "ida"})
	assert.NoError(err)
	assert.True(called)
}
