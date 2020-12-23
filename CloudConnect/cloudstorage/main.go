package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/infinimesh/plugins/CloudConnect/csvprocessor"

	"cloud.google.com/go/storage"
)

const (
	envGcpProjectID  = "GCP_PROJECT_ID"
	envGcpBucketName = "GCP_BUCKET_NAME"
)

func main() {
	bucketName := os.Getenv(envGcpBucketName)
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_ = c.Bucket(bucketName).Create(ctx, os.Getenv(envGcpProjectID), nil)
	csvprocessor.WalkLoop(csvprocessor.WalkFunc(func(f *os.File) error {
		w := c.Bucket(bucketName).Object(f.Name()).NewWriter(context.Background())
		defer w.Close()
		_, err := io.Copy(w, f)
		return err
	}))
}
