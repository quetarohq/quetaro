package testutil

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var (
	AwsRegion          = "us-east-1"
	AwsEndpointUrl     = "http://localhost:4566"
	SqsWaitTimeSeconds = 2
)

func init() {
	if region := os.Getenv("TEST_AWS_REGION"); region != "" {
		AwsRegion = region
	}
	if url := os.Getenv("TEST_AWS_ENDPOINT_URL"); url != "" {
		AwsEndpointUrl = url
	}
	if sec := os.Getenv("TEST_SQS_WAIT_TIME_SEC"); sec != "" {
		n, err := strconv.Atoi(sec)
		if err != nil {
			panic(err)
		}
		SqsWaitTimeSeconds = n
	}
}

func NewAwsConfig(t *testing.T) aws.Config {
	t.Helper()
	optFns := []func(*config.LoadOptions) error{}
	optFns = append(optFns, config.WithRegion(AwsRegion))
	awsCfg, err := config.LoadDefaultConfig(context.Background(), optFns...)

	if err != nil {
		panic(err)
	}

	return awsCfg
}

func NewSQSClient(t *testing.T) *sqs.Client {
	t.Helper()
	awsCfg := NewAwsConfig(t)
	return sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		o.BaseEndpoint = &AwsEndpointUrl
	})
}

func SendMessage(t *testing.T, client *sqs.Client, queueName string, funcName string, body string) {
	t.Helper()
	attrs := map[string]types.MessageAttributeValue{}

	if funcName != "" {
		attrs["FunctionName"] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(funcName),
		}
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(queueName),
		MessageAttributes: attrs,
		MessageBody:       aws.String(body),
	}

	_, err := client.SendMessage(context.Background(), input)

	if err != nil {
		panic(err)
	}
}

func PurgeQueue(t *testing.T, client *sqs.Client, queueName string) {
	t.Helper()

	input := &sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueName),
	}

	_, err := client.PurgeQueue(context.Background(), input)

	if err != nil {
		panic(err)
	}
}

func ReceiveMessages(t *testing.T, client *sqs.Client, queueName string) []types.Message {
	t.Helper()

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueName),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     int32(SqsWaitTimeSeconds),
	}

	output, err := client.ReceiveMessage(context.Background(), input)

	if err != nil {
		panic(err)
	}

	return output.Messages
}

func NewLambdClient(t *testing.T) *lambda.Client {
	t.Helper()
	awsCfg := NewAwsConfig(t)
	return lambda.NewFromConfig(awsCfg)
}
