package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

const (
	envReadDir       = "READ_DIR"
	envGcpProjectID  = "GCP_PROJECT_ID"
	envGcpBucketName = "GCP_BUCKET_NAME"
)

func main() {
	readDir := os.Getenv(envReadDir)
	bucketName := os.Getenv(envGcpBucketName)
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_ = c.Bucket(bucketName).Create(ctx, os.Getenv(envGcpProjectID), nil)

	h := &cloudStorageHandler{
		svc:        c,
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

type cloudStorageHandler struct {
	svc        *storage.Client
	bucketName string
}

func (h *cloudStorageHandler) walkFunc(path string, info os.FileInfo, err error) error {
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
		w := h.svc.Bucket(h.bucketName).Object(path).NewWriter(context.Background())
		defer w.Close()
		_, err = io.Copy(w, f)
		if err != nil {
			log.Printf("failed to write object to cloud storage: %v\n", err)
			return err
		} else {
			log.Printf("successfully put %s to cloud storage\n", path)
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
