package blob

import (
	"NimbusDb/configurations"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the MinIO client and provides blob storage operations.
type Client struct {
	minioClient minioClientInterface
	config      *configurations.Config
}

// NewClient creates a new MinIO client with the provided configuration.
//
// params:
//   - ctx: Context for the operation
//   - cfg: Configuration containing MinIO endpoint, credentials, and bucket name
//
// return:
//   - *Client: A new blob client instance
//   - error: An error if the client could not be initialized
func NewClient(ctx context.Context, cfg *configurations.Config) (*Client, error) {
	if cfg.Blob.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if cfg.Blob.AccessKeyID == "" {
		return nil, fmt.Errorf("access key ID is required")
	}
	if cfg.Blob.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret access key is required")
	}

	minioClient, err := minio.New(cfg.Blob.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Blob.AccessKeyID, cfg.Blob.SecretAccessKey, ""),
		Secure: cfg.Blob.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Test connection with timeout
	testCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = minioClient.ListBuckets(testCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	return &Client{
		minioClient: newMinioClientAdapter(minioClient),
		config:      cfg,
	}, nil
}

// NewClientWithInterface creates a new Client with a custom MinIO client interface.
// This is primarily used for testing with mock implementations.
//
// params:
//   - minioClient: An implementation of minioClientInterface (can be a mock)
//   - cfg: Optional configuration. If nil, operations requiring config (like CreateBucket) will fail
//
// return:
//   - *Client: A new blob client instance
func NewClientWithInterface(minioClient minioClientInterface, cfg *configurations.Config) *Client {
	return &Client{
		minioClient: minioClient,
		config:      cfg,
	}
}
