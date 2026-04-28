package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type Storage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) error
	GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}

type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	creds := credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSS3Region),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.S3UsePathStyle
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = &cfg.S3Endpoint
		}
	})

	presignClient := s3.NewPresignClient(client)

	return &S3Storage{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.S3Bucket,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	}
	_, err := s.client.PutObject(ctx, input)
	return err
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}
	_, err := s.client.HeadObject(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			var notFound *s3types.NotFound
			if errors.As(err, &notFound) || apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "404" {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}
	res, err := s.presignClient.PresignGetObject(ctx, input, func(po *s3.PresignOptions) {
		po.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}
	_, err := s.client.DeleteObject(ctx, input)
	return err
}
