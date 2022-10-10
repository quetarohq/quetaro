package sqsmsg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro/sqsmsg"
)

func Test_DecodeId(t *testing.T) {
	assert := assert.New(t)
	body := `
		{
			"version": "1.0",
			"timestamp": "2022-12-07T16:47:47.153Z",
			"requestContext": {
				"requestId": "808e578c-cbaa-409b-a953-0a050354c76b",
				"functionArn": "arn:aws:lambda:us-east-1:000000000000:function:qtr-job-test",
				"condition": "Success",
				"approximateInvokeCount": 1
			},
			"requestPayload": {
				"_id": "013eb466-184c-43e6-b0c2-6667d5cf3b47"
			},
			"responseContext": {
				"statusCode": 200,
				"executedVersion": "$LATEST"
			},
			"responsePayload": "undefined"
		}
`
	id, err := sqsmsg.DecodeId(body, "_id")
	assert.NoError(err)
	assert.Equal("013eb466-184c-43e6-b0c2-6667d5cf3b47", id)
}
