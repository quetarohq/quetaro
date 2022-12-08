package quetaro_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro"
	"github.com/winebarrel/quetaro/internal/testutil"
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
	assert.NoError(err)
}
