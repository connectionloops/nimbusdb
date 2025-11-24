package blob

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// minioClientInterface defines the interface for MinIO client operations.
// This interface allows us to mock MinIO behavior in unit tests.
type minioClientInterface interface {
	// ListBuckets lists all buckets.
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)

	// GetObject retrieves an object from a bucket.
	// Returns an io.ReadCloser that should be closed after use.
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error)

	// PutObject uploads an object to a bucket.
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)

	// BucketExists checks if a bucket exists.
	BucketExists(ctx context.Context, bucketName string) (bool, error)

	// MakeBucket creates a new bucket.
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error

	// EnableVersioning enables versioning on a bucket.
	EnableVersioning(ctx context.Context, bucketName string) error

	// GetBucketVersioning gets the versioning configuration of a bucket.
	GetBucketVersioning(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error)

	// RemoveObject removes an object from a bucket.
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error

	// RemoveBucket removes a bucket.
	RemoveBucket(ctx context.Context, bucketName string) error

	// SetBucketLifecycle sets the lifecycle configuration for a bucket.
	SetBucketLifecycle(ctx context.Context, bucketName string, config *lifecycle.Configuration) error

	// StatObject retrieves object metadata without reading the object.
	StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
}
