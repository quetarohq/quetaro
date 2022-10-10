package quetaro_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro"
	"github.com/winebarrel/quetaro/internal/testutil"
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
		select {
		case <-time.After(1 * time.Second):
			cancel()
		}
	}()

	err := intakeInvoke.Start(ctx)
	assert.NoError(err)
}
