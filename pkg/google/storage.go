package google

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

// uploadFile uploads an object (file) to a bucket.
func UploadFile(bucket, object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(object)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Upload file to bucket.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("gs://%s/%s", bucket, object), nil
}

// deleteFile deletes an object (file) from a bucket.
func DeleteFile(bucket, object string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Bucket(bucket).Object(object).Delete(ctx)
}
