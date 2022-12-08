package quetaro

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/quetarohq/quetaro/awsutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type IntakeInvokeOpts struct {
	QueueName      string
	ConnConfig     *pgx.ConnConfig
	NAgents        int
	Interval       time.Duration
	ErrInterval    time.Duration
	AWSRegion      string
	AWSEndpointUrl string
}

func (opts *IntakeInvokeOpts) MarshalZerologObject(e *zerolog.Event) {
	e.Str("queue", opts.QueueName).
		Str("dsn", opts.ConnConfig.ConnString()).
		Str("aws_region", opts.AWSRegion).
		Str("aws_endpoint_url", opts.AWSEndpointUrl)
}

type IntakeInvoke struct {
	*IntakeInvokeOpts
	lambda awsutil.LambdaAPI
}

func NewIntakeInvoke(opts *IntakeInvokeOpts) (*IntakeInvoke, error) {
	cfg, err := awsutil.LoadDefaultConfig(opts.AWSRegion, opts.AWSEndpointUrl)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	intakeInvoke := &IntakeInvoke{
		IntakeInvokeOpts: opts,
		lambda:           lambda.NewFromConfig(cfg),
	}

	return intakeInvoke, nil
}

func (intakeInvoke *IntakeInvoke) Start(ctx context.Context) error {
	logger := log.Ctx(ctx).With().Str("queue_name", intakeInvoke.QueueName).Logger()
	ctx = logger.WithContext(ctx)
	logger.Info().Msg("start intake-invoke")
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < intakeInvoke.NAgents; i++ {
		invokeAgent := newIntakeInvokeAgent(intakeInvoke)

		eg.Go(func() error {
			return invokeAgent.run(ctx)
		})
	}

	err := eg.Wait()
	logger.Info().Msg("shutdown intake-invoke")

	if err != nil {
		return errors.Wrap(err, "error in agent")
	}

	return nil
}
