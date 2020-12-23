package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/infinimesh/plugins/CloudConnect/csvprocessor"
)

const (
	envAzAccountName   = "AZ_ACCOUNT_NAME"
	envAzAccountKey    = "AZ_ACCOUNT_KEY"
	envAzContainerName = "AZ_CONTAINER_NAME"
)

func main() {
	accountName := os.Getenv(envAzAccountName)
	accountKey := os.Getenv(envAzAccountKey)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	serviceURL := azblob.NewServiceURL(*u, p)
	ctx := context.Background()
	containerURL := serviceURL.NewContainerURL(os.Getenv(envAzContainerName))
	_, _ = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	csvprocessor.WalkLoop(csvprocessor.WalkFunc(func(f *os.File) error {
		_, err = containerURL.
			NewBlockBlobURL(f.Name()).
			Upload(context.Background(), f, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil, azblob.ClientProvidedKeyOptions{})
		return err
	}))
}
