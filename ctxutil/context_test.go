package ctxutil_test

import (
	"context"
	"testing"

	"github.com/quetarohq/quetaro/ctxutil"
	"github.com/stretchr/testify/assert"
)

func Test_WithoutCancel(t *testing.T) {
	assert := assert.New(t)
	parent, cancel := context.WithCancel(context.Background())
	ctx := ctxutil.WithoutCancel(parent)
	cancel()
	assert.Nil(ctx.Done())
	assert.Equal(struct{}{}, <-parent.Done())
}
