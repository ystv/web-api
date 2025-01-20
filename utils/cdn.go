package utils

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CDNConfig represents a configuration to connect to a CDN / S3 instance
type CDNConfig struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// NewCDN Initialise connection to CDN
func NewCDN(config CDNConfig) (*s3.Client, error) {
	s3Config, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(config.Region),
		awsConfig.WithBaseEndpoint(config.Endpoint),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	cdn := s3.NewFromConfig(s3Config)
	return cdn, nil
}
