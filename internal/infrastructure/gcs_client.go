package infrastructure

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// NewGCSClient creates a new Google Cloud Storage client
func NewGCSClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return client, nil
}

// CloseGCSClient closes the GCS client
func CloseGCSClient(client *storage.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}
