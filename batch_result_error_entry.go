package quetaro

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog"
)

type BatchResultErrorEntry types.BatchResultErrorEntry

func (bree BatchResultErrorEntry) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("batch_result_error_entry_code", aws.ToString(bree.Code)).
		Str("batch_result_error_entry_id", aws.ToString(bree.Id)).
		Str("batch_result_error_entry_message", aws.ToString(bree.Message))
}

func FromBatchResultErrorEntry(bree types.BatchResultErrorEntry) *BatchResultErrorEntry {
	obj := BatchResultErrorEntry(bree)
	return &obj
}
