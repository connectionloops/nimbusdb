package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/minio/minio-go/v7"
)

var (
	// bucketNameRegex validates bucket names according to S3/MinIO naming rules:
	// - 3-63 characters
	// - Lowercase letters, numbers, dots, and hyphens only
	// - Must start and end with a letter or number
	// - No consecutive dots
	bucketNameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$`)
)

// validateBucketName validates a bucket name according to S3/MinIO naming rules.
// only used during create bucket operation
func validateBucketName(bucketName string) error {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("bucket name must be between 3 and 63 characters, got length %d", len(bucketName))
	}

	if !bucketNameRegex.MatchString(bucketName) {
		return fmt.Errorf("bucket name contains invalid characters or format: %s", bucketName)
	}

	// Check for consecutive dots
	if strings.Contains(bucketName, "..") {
		return fmt.Errorf("bucket name cannot contain consecutive dots: %s", bucketName)
	}

	return nil
}

// ReadFile reads a file from MinIO and returns its contents as a byte array.
//
// params:
//   - ctx: Context for the operation
//   - bucketName: The name of the bucket to read from
//   - fileName: The name of the file to read
//
// return:
//   - []byte: The file contents
//   - error: An error if the file could not be read
func (c *Client) ReadFile(ctx context.Context, bucketName, fileName string) ([]byte, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}
	if fileName == "" {
		return nil, fmt.Errorf("file name cannot be empty")
	}

	object, err := c.minioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", fileName, err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", fileName, err)
	}

	return data, nil
}

// WriteFile writes a byte array to a file in MinIO.
//
// params:
//   - ctx: Context for the operation
//   - bucketName: The name of the bucket to write to
//   - fileName: The name of the file to write
//   - data: The data to write
//
// return:
//   - error: An error if the file could not be written
func (c *Client) WriteFile(ctx context.Context, bucketName, fileName string, data []byte) error {
	if bucketName == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	_, err := c.minioClient.PutObject(ctx, bucketName, fileName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("failed to put object %s: %w", fileName, err)
	}

	return nil
}

// CreateBucket creates a new bucket in MinIO with versioning enabled.
//
// params:
//   - ctx: Context for the operation
//   - bucketName: The name of the bucket to create
//
// return:
//   - error: An error if the bucket could not be created or versioning could not be enabled
func (c *Client) CreateBucket(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}
	if err := validateBucketName(bucketName); err != nil {
		return err
	}

	// Check if bucket already exists
	exists, err := c.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		// Create the bucket
		err = c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	// Enable versioning on the bucket
	err = c.minioClient.EnableVersioning(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to enable versioning on bucket %s: %w", bucketName, err)
	}

	return nil
}
