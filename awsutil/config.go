package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadDefaultConfig(region string, endpointUrl string) (aws.Config, error) {
	optFns := []func(*config.LoadOptions) error{}

	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	if endpointUrl != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpointUrl,
				SigningRegion: region,
			}, nil
		})

		optFns = append(optFns, config.WithEndpointResolverWithOptions(customResolver))
	}

	return config.LoadDefaultConfig(context.Background(), optFns...)
}
