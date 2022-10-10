package quetaro

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/winebarrel/quetaro/awsutil"
	"golang.org/x/sync/errgroup"
)

type OutletFailureOpts struct {
	QueueName      string
	ConnConfig     *pgx.ConnConfig
	NAgents        int
	Interval       time.Duration
	ErrInterval    time.Duration
	MaxRecvNum     int
	AWSRegion      string
	AWSEndpointUrl string
}

func (opts *OutletFailureOpts) MarshalZerologObject(e *zerolog.Event) {
	e.Str("queue", opts.QueueName).
		Str("dsn", opts.ConnConfig.ConnString()).
		Str("aws_region", opts.AWSRegion).
		Str("aws_endpoint_url", opts.AWSEndpointUrl)
}

type OutletFailure struct {
	*OutletFailureOpts
	AwsCfg   aws.Config
	QueueUrl string
}

func NewOutletFailure(opts *OutletFailureOpts) (*OutletFailure, error) {
	cfg, err := awsutil.LoadDefaultConfig(opts.AWSRegion, opts.AWSEndpointUrl)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	// SQS client is created by each agent.
	client := sqs.NewFromConfig(cfg)

	// get the queue URL outside the agent.
	output, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(opts.QueueName),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get queue URL")
	}

	outletFailure := &OutletFailure{
		OutletFailureOpts: opts,
		AwsCfg:            cfg,
		QueueUrl:          aws.ToString(output.QueueUrl),
	}

	return outletFailure, nil
}

func (outletFailure *OutletFailure) Start(ctx context.Context) error {
	logger := log.Ctx(ctx).With().Str("queue_name", outletFailure.QueueName).Logger()
	ctx = logger.WithContext(ctx)
	logger.Info().Msg("start outlet-failure")
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < outletFailure.NAgents; i++ {
		failureAgent := newOutletFailureAgent(outletFailure)

		eg.Go(func() error {
			return failureAgent.run(ctx)
		})
	}

	err := eg.Wait()
	logger.Info().Msg("shutdown outlet-failure")

	if err != nil {
		return errors.Wrap(err, "error in agent")
	}

	return nil
}
