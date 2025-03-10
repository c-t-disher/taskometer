package services

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"io"
	"log"
)

// TODO this should be an env var
const bucketName = "taskometer-bucket"

type settings struct {
	s3Client   *s3.Client
	bucketName string
}

var s settings

func InitS3Client() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("taskometer-user"))
	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	s = settings{
		s3Client:   s3Client,
		bucketName: bucketName,
	}

	exists, err := s.bucketExists(context.Background(), s.bucketName)
	if !exists {
		log.Fatal("Task bucket does not exist!")
	}

}

func getS3Client() *s3.Client {
	return s.s3Client
}

// BucketExists checks whether a bucket exists in the current account.
func (s settings) bucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := getS3Client().HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you own it.", bucketName)
	}

	return exists, err
}

func PutObj(ctx context.Context, key string, body io.Reader) error {
	_, err := getS3Client().PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   body,
	})

	return err
}

func GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	res, err := getS3Client().GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error) {
	return getS3Client().ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix + "/"),
	})
}

func DeleteObj(ctx context.Context, key string) error {
	_, err := getS3Client().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}
