package awsutil

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/pkg/errors"
)

func InvokeFunction(ctx context.Context, client LambdaAPI, funcName string, payload0 string, extra map[string]string) error {
	m := map[string]any{}
	err := json.Unmarshal([]byte(payload0), &m)

	if err != nil {
		// must not happen
		panic(err)
	}

	for k, v := range extra {
		m[k] = v
	}

	payload, err := json.Marshal(m)

	if err != nil {
		// must not happen
		panic(err)
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(funcName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	}

	_, err = client.Invoke(ctx, input)

	if err != nil {
		return errors.Wrap(err, "Lambda Invoke error")
	}

	return nil
}

type LambdaAPI interface {
	Invoke(context.Context, *lambda.InvokeInput, ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}
