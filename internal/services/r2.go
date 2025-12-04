package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pavelc4/ytdpl-api-go/config"
)

type R2Service struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewR2Service(cfg config.R2Config) (*R2Service, error) {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("R2 credentials not configured")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: cfg.Endpoint,
		}, nil
	})

	awsCfg := aws.Config{
		Region:                      "auto",
		Credentials:                 credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		EndpointResolverWithOptions: r2Resolver,
	}

	client := s3.NewFromConfig(awsCfg)

	return &R2Service{
		client:    client,
		bucket:    cfg.BucketName,
		publicURL: cfg.PublicURL,
	}, nil
}
func (r *R2Service) UploadFile(ctx context.Context, localPath, objectKey string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", r.publicURL, objectKey)
	return publicURL, nil
}

func (r *R2Service) DeleteFile(ctx context.Context, objectKey string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(objectKey),
	})
	return err
}

func (r *R2Service) CleanupOldFiles(ctx context.Context, retentionDays int) error {
	log.Printf(" Starting cleanup of files older than %d days...", retentionDays)

	prefixes := []string{"vidioe/", "audio/"}
	deletedCount := 0
	errorsCount := 0
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	for _, prefix := range prefixes {
		paginator := s3.NewListObjectsV2Paginator(r.client, &s3.ListObjectsV2Input{
			Bucket: aws.String(r.bucket),
			Prefix: aws.String(prefix),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				log.Printf("Failed to list objects for prefix %s: %v", prefix, err)
				errorsCount++
				continue
			}

			for _, obj := range page.Contents {
				if obj.LastModified.Before(cutoff) {
					err := r.DeleteFile(ctx, *obj.Key)
					if err != nil {
						log.Printf(" Failed to delete %s: %v", *obj.Key, err)
						errorsCount++
					} else {
						log.Printf(" Deleted old file: %s (Last modified: %s)", *obj.Key, obj.LastModified.Format(time.RFC3339))
						deletedCount++
					}
				}
			}
		}
	}

	log.Printf(" Cleanup completed. Deleted: %d, Errors: %d", deletedCount, errorsCount)
	return nil
}

func (r *R2Service) DownloadFile(ctx context.Context, objectKey, destPath string) error {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from R2: %w", err)
	}
	defer result.Body.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, result.Body)
	if err != nil {
		return fmt.Errorf("failed to write to local file: %w", err)
	}

	return nil
}
