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

func Test_IntakePull_Start(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	intakePull := &quetaro.IntakePull{
		IntakePullOpts: &quetaro.IntakePullOpts{
			QueueName:  "qtr-intake-test",
			ConnConfig: res.ConnCfg,
			NAgents:    3,
		},
		AwsCfg:   testutil.NewAwsConfig(t),
		QueueUrl: "qtr-intake-test",
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := intakePull.Start(ctx)
	assert.True(err == nil || errors.Is(err, context.Canceled))
}
