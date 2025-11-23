package blob

import (
	"context"
	"errors"
	"testing"

	"NimbusDb/configurations"
)

func TestNewClient_ValidConfig(t *testing.T) {
	// This test requires a running MinIO instance
	// For integration tests, you would set up a test MinIO server
	// For unit tests, you might want to use mocks instead
	cfg := &configurations.Config{
		Blob: configurations.BlobConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
		},
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err != nil {
		// Skip test if MinIO is not available
		t.Skipf("MinIO not available, skipping test: %v", err)
		return
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClient_EmptyEndpoint(t *testing.T) {
	cfg := &configurations.Config{
		Blob: configurations.BlobConfig{
			Endpoint:        "",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
		},
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err == nil {
		t.Error("NewClient() should have failed with empty endpoint")
	}
	if client != nil {
		t.Error("NewClient() should return nil client on error")
	}
}

func TestNewClient_EmptyAccessKeyID(t *testing.T) {
	cfg := &configurations.Config{
		Blob: configurations.BlobConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
		},
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err == nil {
		t.Error("NewClient() should have failed with empty access key ID")
	}
	if client != nil {
		t.Error("NewClient() should return nil client on error")
	}
}

func TestNewClient_EmptySecretAccessKey(t *testing.T) {
	cfg := &configurations.Config{
		Blob: configurations.BlobConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "",
			UseSSL:          false,
		},
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err == nil {
		t.Error("NewClient() should have failed with empty secret access key")
	}
	if client != nil {
		t.Error("NewClient() should return nil client on error")
	}
}

func TestNewClientWithInterface(t *testing.T) {
	mockClient := newMockMinioClient()
	testConfig := &configurations.Config{
		Blob: configurations.BlobConfig{
			DeleteMarkerCleanupDelayDays:      1,
			NonCurrentVersionCleanupDelayDays: 1,
		},
	}
	client := NewClientWithInterface(mockClient, testConfig)

	if client == nil {
		t.Fatal("NewClientWithInterface() returned nil client")
	}

	if client.minioClient == nil {
		t.Fatal("NewClientWithInterface() returned client with nil minioClient")
	}

	if client.config == nil {
		t.Fatal("NewClientWithInterface() returned client with nil config")
	}
}

func TestNewClient_ConnectionFailure(t *testing.T) {
	cfg := &configurations.Config{
		Blob: configurations.BlobConfig{
			Endpoint:        "invalid-endpoint:9999",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
		},
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err == nil {
		t.Error("NewClient() should have failed with invalid endpoint")
	}
	if client != nil {
		t.Error("NewClient() should return nil client on error")
	}
}

func TestNewClient_ListBucketsError(t *testing.T) {
	// Test that ListBuckets error is properly handled
	mockClient := newMockMinioClient()
	mockClient.setListBucketsError(errors.New("connection failed"))

	// This test verifies that the mock can simulate connection failures
	ctx := context.Background()
	_, err := mockClient.ListBuckets(ctx)
	if err == nil {
		t.Error("ListBuckets() should have failed with mock error")
	}
}
