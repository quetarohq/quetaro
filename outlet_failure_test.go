package quetaro_test

import (
	"context"
	"testing"
	"time"

	"github.com/quetarohq/quetaro"
	"github.com/quetarohq/quetaro/internal/testutil"
	"github.com/stretchr/testify/assert"
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
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := outletFailure.Start(ctx)
	assert.NoError(err)
}
