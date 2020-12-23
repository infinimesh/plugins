package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/infinimesh/plugins/CloudConnect/csvprocessor"
)

const envS3BucketName = "AWS_S3BUCKETNAME"

func main() {
	bucketName := os.Getenv(envS3BucketName)
	svc := s3.New(session.New())
	_, _ = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	csvprocessor.WalkLoop(csvprocessor.WalkFunc(func(f *os.File) error {
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(f.Name()),
			Body:   f,
		})
		return err
	}))
}
