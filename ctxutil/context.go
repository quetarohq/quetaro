package ctxutil

import (
	"context"
	"time"
)

type withoutCancelCtx struct {
	parent context.Context
}

func (*withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (*withoutCancelCtx) Err() error {
	return nil
}

func (ctx *withoutCancelCtx) Value(key any) any {
	return ctx.parent.Value(key)
}

func WithoutCancel(parent context.Context) context.Context {
	return &withoutCancelCtx{parent}
}
