package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"personae-fasti/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ImageStorage struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
}

func NewImageStorage(config *configs.ImageStorageConfig) *ImageStorage {
	if config == nil || config.Endpoint == "" || config.Bucket == "" || config.AccessKey == "" || config.SecretKey == "" || config.PublicBaseURL == "" {
		return &ImageStorage{}
	}
	region := config.Region
	if region == "" {
		region = "garage"
	}
	awsConfig := aws.Config{
		Region:                     region,
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			config.AccessKey,
			config.SecretKey,
			"",
		)),
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(config.Endpoint)
		options.UsePathStyle = config.ForcePathStyle
	})
	return &ImageStorage{
		client:        client,
		bucket:        config.Bucket,
		publicBaseURL: strings.TrimRight(config.PublicBaseURL, "/"),
	}
}

func (s *ImageStorage) Configured() bool {
	return s != nil && s.client != nil
}

func (s *ImageStorage) Put(ctx context.Context, key, contentType string, data []byte) error {
	if !s.Configured() {
		return errors.New("image storage is not configured")
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(data),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	return err
}

func (s *ImageStorage) Delete(ctx context.Context, key string) error {
	if !s.Configured() || key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *ImageStorage) PublicURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicBaseURL, strings.TrimLeft(key, "/"))
}
