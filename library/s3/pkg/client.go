package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func provideS3Client() (*s3.Client, error) {
	conf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(conf), nil
}

func provideS3PresignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}
