package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"time"
)

const (
	configSectionName = "s3"
)

type linkGeneratorConfig struct {
	Bucket          string        `yaml:"bucket"`
	DownloadLinkTTL time.Duration `yaml:"download_link_ttl"`
	UploadLinkTTL   time.Duration `yaml:"upload_link_ttl"`
}

type LinkGenerator interface {
	GenerateDownloadLink(ctx context.Context, key string) (link string, method string, err error)
	GenerateUploadLink(ctx context.Context, key string) (link string, method string, err error)
}

func provideLinkGenerator(
	s3Client *s3.PresignClient,
	config linkGeneratorConfig,
) LinkGenerator {
	return &linkGenerator{
		client: s3Client,
		config: config,
	}
}

type linkGenerator struct {
	client *s3.PresignClient
	config linkGeneratorConfig
}

func (gen *linkGenerator) GenerateDownloadLink(ctx context.Context, key string) (link string, method string, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := gen.client.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Key:    aws.String(key),
			Bucket: aws.String(gen.config.Bucket),
		},
		s3.WithPresignExpires(gen.config.DownloadLinkTTL),
	)

	if err != nil {
		return "", "", errors.WithStack(err)
	}

	return req.URL, req.Method, nil
}

func (gen *linkGenerator) GenerateUploadLink(ctx context.Context, key string) (link string, method string, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := gen.client.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Key:    aws.String(key),
			Bucket: aws.String(gen.config.Bucket),
		},
		s3.WithPresignExpires(gen.config.UploadLinkTTL),
	)

	if err != nil {
		return "", "", err
	}

	return req.URL, req.Method, errors.WithStack(err)
}
