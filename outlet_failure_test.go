package quetaro_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/quetaro"
	"github.com/winebarrel/quetaro/internal/testutil"
)

func Test_OutletFailure_Start(t *testing.T) {
	res := testutil.SetupAgent(t)
	defer res.Cleanup(t)
	assert := assert.New(t)

	outletFailure := &quetaro.OutletFailure{
		OutletFailureOpts: &quetaro.OutletFailureOpts{
			QueueName:  "qtr-outlet-failure-test",
			ConnConfig: res.ConnCfg,
			NAgents:    3,
		},
		AwsCfg:   testutil.NewAwsConfig(t),
		QueueUrl: "qtr-outlet-failure-test",
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-time.After(1 * time.Second):
			cancel()
		}
	}()

	err := outletFailure.Start(ctx)
	assert.NoError(err)
}
