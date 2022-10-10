package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
)

func ReceiveMessages(ctx context.Context, client SQSAPI, queueUrl string, maxRecvNum int, attrNames ...string) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueUrl),
		MaxNumberOfMessages:   int32(maxRecvNum),
		MessageAttributeNames: attrNames,
	}

	output, err := client.ReceiveMessage(ctx, input)

	if err != nil {
		return nil, errors.Wrap(err, "SQS ReceiveMessage error")
	}

	return output.Messages, nil
}

func DeleteMessages(ctx context.Context, client SQSAPI, queueUrl string, msgs []types.Message) (error, []types.BatchResultErrorEntry) {
	input := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(queueUrl),
		Entries:  make([]types.DeleteMessageBatchRequestEntry, 0, len(msgs)),
	}

	for _, m := range msgs {
		input.Entries = append(input.Entries, types.DeleteMessageBatchRequestEntry{
			Id:            m.MessageId,
			ReceiptHandle: m.ReceiptHandle,
		})
	}

	resp, err := client.DeleteMessageBatch(ctx, input)

	if err != nil {
		return errors.Wrap(err, "SQS DeleteMessageBatch error"), resp.Failed
	}

	return nil, nil
}

type SQSAPI interface {
	ReceiveMessage(context.Context, *sqs.ReceiveMessageInput, ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessageBatch(context.Context, *sqs.DeleteMessageBatchInput, ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error)
}
