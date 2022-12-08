package dbutil_test

import (
	"context"
	"testing"

	"github.com/quetarohq/quetaro/dbutil"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_Connect(t *testing.T) {
	assert := assert.New(t)

	connCfg := testutil.NewConnConfig(t)
	conn, err := dbutil.Connect(context.Background(), connCfg)
	assert.NoError(err)
	defer conn.Close(context.Background())

	var n int
	err = conn.QueryRow(context.Background(), "select 1").Scan(&n)
	assert.NoError(err)
	assert.Equal(1, n)
}
