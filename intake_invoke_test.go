package quetaro_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/quetarohq/quetaro"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_IntakeInvoke_Start(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	intakeInvoke := &quetaro.IntakeInvoke{
		IntakeInvokeOpts: &quetaro.IntakeInvokeOpts{
			QueueName:  "qtr-intake-test",
			ConnConfig: res.ConnCfg,
			NAgents:    3,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := intakeInvoke.Start(ctx)
	assert.True(err == nil || errors.Is(err, context.Canceled))
}
