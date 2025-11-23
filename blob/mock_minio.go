package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// mockMinioClient is a mock implementation of minioClientInterface for testing.
type mockMinioClient struct {
	mu                  sync.RWMutex
	buckets             map[string]bool
	objects             map[string]map[string][]byte            // bucket -> object -> latest data
	objectVersions      map[string]map[string]map[string][]byte // bucket -> object -> versionID -> data
	latestVersions      map[string]map[string]string            // bucket -> object -> latest versionID
	versioning          map[string]bool                         // bucket -> versioning enabled
	versionCounter      atomic.Int64                            // counter for generating version IDs
	listBucketsErr      error
	getObjectErr        map[string]error                    // bucket/object -> error
	putObjectErr        map[string]error                    // bucket/object -> error
	bucketExistsErr     map[string]error                    // bucket -> error
	makeBucketErr       map[string]error                    // bucket -> error
	enableVersioningErr map[string]error                    // bucket -> error
	removeObjectErr     map[string]error                    // bucket/object -> error
	removeBucketErr     map[string]error                    // bucket -> error
	setLifecycleErr     map[string]error                    // bucket -> error
	lifecycleConfigs    map[string]*lifecycle.Configuration // bucket -> lifecycle config
}

// newMockMinioClient creates a new mock MinIO client.
func newMockMinioClient() *mockMinioClient {
	return &mockMinioClient{
		buckets:             make(map[string]bool),
		objects:             make(map[string]map[string][]byte),
		objectVersions:      make(map[string]map[string]map[string][]byte),
		latestVersions:      make(map[string]map[string]string),
		versioning:          make(map[string]bool),
		getObjectErr:        make(map[string]error),
		putObjectErr:        make(map[string]error),
		bucketExistsErr:     make(map[string]error),
		makeBucketErr:       make(map[string]error),
		enableVersioningErr: make(map[string]error),
		removeObjectErr:     make(map[string]error),
		removeBucketErr:     make(map[string]error),
		setLifecycleErr:     make(map[string]error),
		lifecycleConfigs:    make(map[string]*lifecycle.Configuration),
	}
}

// ListBuckets lists all buckets.
func (m *mockMinioClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	if m.listBucketsErr != nil {
		return nil, m.listBucketsErr
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	buckets := make([]minio.BucketInfo, 0, len(m.buckets))
	for bucketName := range m.buckets {
		buckets = append(buckets, minio.BucketInfo{
			Name: bucketName,
		})
	}

	return buckets, nil
}

// GetObject retrieves an object from a bucket.
func (m *mockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	if err, ok := m.getObjectErr[key]; ok {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	bucket, exists := m.objects[bucketName]
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucketName)
	}

	var data []byte
	var found bool

	// If version ID is specified, read that specific version
	if opts.VersionID != "" {
		if m.objectVersions[bucketName] != nil && m.objectVersions[bucketName][objectName] != nil {
			data, found = m.objectVersions[bucketName][objectName][opts.VersionID]
		}
		if !found {
			return nil, fmt.Errorf("version %s does not exist for object %s in bucket %s", opts.VersionID, objectName, bucketName)
		}
	} else {
		// No version specified, read latest version
		data, found = bucket[objectName]
		if !found {
			return nil, fmt.Errorf("object %s does not exist in bucket %s", objectName, bucketName)
		}
	}

	// Create a mock object that implements io.ReadCloser
	// We use bytes.NewReader wrapped in io.NopCloser to create a ReadCloser
	return io.NopCloser(bytes.NewReader(data)), nil
}

// PutObject uploads an object to a bucket.
func (m *mockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	if err, ok := m.putObjectErr[key]; ok {
		return minio.UploadInfo{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.buckets[bucketName] {
		return minio.UploadInfo{}, fmt.Errorf("bucket %s does not exist", bucketName)
	}

	if m.objects[bucketName] == nil {
		m.objects[bucketName] = make(map[string][]byte)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	m.objects[bucketName][objectName] = data

	// Generate version ID if versioning is enabled
	var versionID string
	if m.versioning[bucketName] {
		versionID = fmt.Sprintf("version-%d", m.versionCounter.Add(1))

		// Initialize version tracking structures if needed
		if m.objectVersions[bucketName] == nil {
			m.objectVersions[bucketName] = make(map[string]map[string][]byte)
		}
		if m.objectVersions[bucketName][objectName] == nil {
			m.objectVersions[bucketName][objectName] = make(map[string][]byte)
		}
		if m.latestVersions[bucketName] == nil {
			m.latestVersions[bucketName] = make(map[string]string)
		}

		// Store the versioned data
		m.objectVersions[bucketName][objectName][versionID] = data
		m.latestVersions[bucketName][objectName] = versionID
	}

	return minio.UploadInfo{
		Bucket:    bucketName,
		Key:       objectName,
		Size:      int64(len(data)),
		VersionID: versionID,
	}, nil
}

// BucketExists checks if a bucket exists.
func (m *mockMinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	if err, ok := m.bucketExistsErr[bucketName]; ok {
		return false, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.buckets[bucketName], nil
}

// MakeBucket creates a new bucket.
func (m *mockMinioClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	if err, ok := m.makeBucketErr[bucketName]; ok {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.buckets[bucketName] {
		return fmt.Errorf("bucket %s already exists", bucketName)
	}

	m.buckets[bucketName] = true
	m.objects[bucketName] = make(map[string][]byte)

	return nil
}

// EnableVersioning enables versioning on a bucket.
func (m *mockMinioClient) EnableVersioning(ctx context.Context, bucketName string) error {
	if err, ok := m.enableVersioningErr[bucketName]; ok {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.buckets[bucketName] {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	m.versioning[bucketName] = true
	return nil
}

// GetBucketVersioning gets the versioning configuration of a bucket.
func (m *mockMinioClient) GetBucketVersioning(ctx context.Context, bucketName string) (minio.BucketVersioningConfiguration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.buckets[bucketName] {
		return minio.BucketVersioningConfiguration{}, fmt.Errorf("bucket %s does not exist", bucketName)
	}

	status := "Disabled"
	if m.versioning[bucketName] {
		status = "Enabled"
	}

	return minio.BucketVersioningConfiguration{
		Status: status,
	}, nil
}

// RemoveObject removes an object from a bucket.
func (m *mockMinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	if err, ok := m.removeObjectErr[key]; ok {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	bucket, exists := m.objects[bucketName]
	if !exists {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	delete(bucket, objectName)
	return nil
}

// RemoveBucket removes a bucket.
func (m *mockMinioClient) RemoveBucket(ctx context.Context, bucketName string) error {
	if err, ok := m.removeBucketErr[bucketName]; ok {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.buckets[bucketName] {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	delete(m.buckets, bucketName)
	delete(m.objects, bucketName)
	delete(m.versioning, bucketName)
	delete(m.lifecycleConfigs, bucketName)
	delete(m.objectVersions, bucketName)
	delete(m.latestVersions, bucketName)

	return nil
}

// SetBucketLifecycle sets the lifecycle configuration for a bucket.
func (m *mockMinioClient) SetBucketLifecycle(ctx context.Context, bucketName string, config *lifecycle.Configuration) error {
	if err, ok := m.setLifecycleErr[bucketName]; ok {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.buckets[bucketName] {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	m.lifecycleConfigs[bucketName] = config
	return nil
}

// Helper methods for test setup

// setListBucketsError sets an error to return from ListBuckets.
func (m *mockMinioClient) setListBucketsError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listBucketsErr = err
}

// setGetObjectError sets an error to return from GetObject for a specific bucket/object.
func (m *mockMinioClient) setGetObjectError(bucketName, objectName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	m.getObjectErr[key] = err
}

// setPutObjectError sets an error to return from PutObject for a specific bucket/object.
func (m *mockMinioClient) setPutObjectError(bucketName, objectName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	m.putObjectErr[key] = err
}

// setBucketExistsError sets an error to return from BucketExists for a specific bucket.
func (m *mockMinioClient) setBucketExistsError(bucketName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bucketExistsErr[bucketName] = err
}

// setMakeBucketError sets an error to return from MakeBucket for a specific bucket.
func (m *mockMinioClient) setMakeBucketError(bucketName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.makeBucketErr[bucketName] = err
}

// setEnableVersioningError sets an error to return from EnableVersioning for a specific bucket.
func (m *mockMinioClient) setEnableVersioningError(bucketName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enableVersioningErr[bucketName] = err
}

// createBucketForTesting creates a bucket for testing purposes.
func (m *mockMinioClient) createBucketForTesting(bucketName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.buckets[bucketName] = true
	m.objects[bucketName] = make(map[string][]byte)
	m.versioning[bucketName] = true // Enable versioning by default for tests
	m.objectVersions[bucketName] = make(map[string]map[string][]byte)
	m.latestVersions[bucketName] = make(map[string]string)
}
