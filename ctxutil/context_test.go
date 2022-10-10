package ctxutil_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro/ctxutil"
)

func Test_WithoutCancel(t *testing.T) {
	assert := assert.New(t)
	parent, cancel := context.WithCancel(context.Background())
	ctx := ctxutil.WithoutCancel(parent)
	cancel()
	assert.Nil(ctx.Done())
	assert.Equal(struct{}{}, <-parent.Done())
}
