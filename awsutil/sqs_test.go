package awsutil_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestReceiveMessage(t *testing.T) {
	assert := assert.New(t)
	var called bool

	mock := &testutil.MockSQS{
		MockReceiveMessage: func(ctx context.Context, rmi *sqs.ReceiveMessageInput, f ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
			called = true
			assert.Equal("queueUrl", aws.ToString(rmi.QueueUrl))
			assert.Equal(int32(3), rmi.MaxNumberOfMessages)
			assert.Equal([]string{"attrName"}, rmi.MessageAttributeNames)
			return &sqs.ReceiveMessageOutput{}, nil
		},
	}

	_, err := awsutil.ReceiveMessages(context.Background(), mock, "queueUrl", 3, "attrName")
	assert.NoError(err)
	assert.True(called)
}

func TestDeleteMessageBatch(t *testing.T) {
	assert := assert.New(t)
	var called bool

	mock := &testutil.MockSQS{
		MockDeleteMessageBatch: func(ctx context.Context, dmbi *sqs.DeleteMessageBatchInput, f ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error) {
			called = true
			assert.Equal("ID", aws.ToString(dmbi.Entries[0].Id))
			return &sqs.DeleteMessageBatchOutput{
				Failed: []types.BatchResultErrorEntry{{Id: aws.String("FailedID")}},
			}, errors.New("DeleteMessageBatch error")
		},
	}

	msgs := []types.Message{{MessageId: aws.String("ID")}}
	failed, err := awsutil.DeleteMessages(context.Background(), mock, "queueUrl", msgs)
	assert.EqualError(err, "SQS DeleteMessageBatch error: DeleteMessageBatch error")
	assert.True(called)
	assert.Equal(1, len(failed))
	assert.Equal("FailedID", aws.ToString(failed[0].Id))
}
