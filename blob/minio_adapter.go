package blob

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// minioClientAdapter adapts a real minio.Client to implement minioClientInterface.
type minioClientAdapter struct {
	client *minio.Client
}

// newMinioClientAdapter creates a new adapter for a real MinIO client.
func newMinioClientAdapter(client *minio.Client) minioClientInterface {
	return &minioClientAdapter{client: client}
}

// ListBuckets lists all buckets.
func (a *minioClientAdapter) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return a.client.ListBuckets(ctx)
}

// GetObject retrieves an object from a bucket.
func (a *minioClientAdapter) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	return a.client.GetObject(ctx, bucketName, objectName, opts)
}

// PutObject uploads an object to a bucket.
func (a *minioClientAdapter) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return a.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

// BucketExists checks if a bucket exists.
func (a *minioClientAdapter) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return a.client.BucketExists(ctx, bucketName)
}

// MakeBucket creates a new bucket.
func (a *minioClientAdapter) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	return a.client.MakeBucket(ctx, bucketName, opts)
}

// EnableVersioning enables versioning on a bucket.
func (a *minioClientAdapter) EnableVersioning(ctx context.Context, bucketName string) error {
	return a.client.EnableVersioning(ctx, bucketName)
}

// GetBucketVersioning gets the versioning configuration of a bucket.
func (a *minioClientAdapter) GetBucketVersioning(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error) {
	return a.client.GetBucketVersioning(ctx, bucketName)
}

// RemoveObject removes an object from a bucket.
func (a *minioClientAdapter) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return a.client.RemoveObject(ctx, bucketName, objectName, opts)
}

// RemoveBucket removes a bucket.
func (a *minioClientAdapter) RemoveBucket(ctx context.Context, bucketName string) error {
	return a.client.RemoveBucket(ctx, bucketName)
}

// SetBucketLifecycle sets the lifecycle configuration for a bucket.
func (a *minioClientAdapter) SetBucketLifecycle(ctx context.Context, bucketName string, config *lifecycle.Configuration) error {
	return a.client.SetBucketLifecycle(ctx, bucketName, config)
}

// StatObject retrieves object metadata without reading the object.
func (a *minioClientAdapter) StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	return a.client.StatObject(ctx, bucketName, objectName, opts)
}
