package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadDefaultConfig(region string) (aws.Config, error) {
	optFns := []func(*config.LoadOptions) error{}

	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	return config.LoadDefaultConfig(context.Background(), optFns...)
}
