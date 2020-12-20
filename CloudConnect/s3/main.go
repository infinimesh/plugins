package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	envReadDir      = "READ_DIR"
	envS3BucketName = "AWS_S3BUCKETNAME"
)

func main() {
	readDir := os.Getenv(envReadDir)
	bucketName := os.Getenv(envS3BucketName)
	svc := s3.New(session.New())
	_, _ = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	h := &s3Handler{
		svc:        svc,
		bucketName: bucketName,
	}
	for range time.Tick(5 * time.Second) {
		log.Printf("walking %s for new files...", readDir)
		err := filepath.Walk(readDir, h.walkFunc)
		if err != nil {
			log.Printf("failed to walk directory: %v\n", err)
		}
	}
}

type s3Handler struct {
	svc        *s3.S3
	bucketName string
}

func (h *s3Handler) walkFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	// only process files that have not been modified recently
	if time.Now().Add(-1 * time.Minute).After(info.ModTime()) {
		log.Println("processing file:", path)
		f, err := os.Open(path)
		if err != nil {
			log.Printf("failed to open %s: %v\n", path, err)
			return err
		}
		_, err = h.svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(h.bucketName),
			Key:    aws.String(path),
			Body:   f,
		})
		if err != nil {
			log.Printf("failed to put object to s3: %v\n", err)
			return err
		} else {
			log.Printf("successfully put %s to s3\n", path)
		}
		err = os.Remove(path)
		if err != nil {
			log.Printf("failed to remove object from path: %v\n", err)
			return err
		} else {
			log.Printf("successfully deleted %s\n", path)
		}
	}
	return nil
}
