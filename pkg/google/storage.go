package google

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	// Get the base name of the file (without directory path).
	baseFilename := filepath.Base(object)

	// Upload file to bucket.
	wc := client.Bucket(bucket).Object(baseFilename).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("gs://%s/%s", bucket, baseFilename), nil
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
