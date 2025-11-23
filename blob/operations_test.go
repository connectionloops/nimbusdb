package blob

import (
	"context"
	"testing"
)

// setupMockClient creates a test client with a mock MinIO implementation.
func setupMockClient(t *testing.T) (*Client, string) {
	mockClient := newMockMinioClient()
	bucketName := "test-bucket"

	// Pre-create the bucket for testing
	mockClient.createBucketForTesting(bucketName)

	client := NewClientWithInterface(mockClient)
	return client, bucketName
}

func TestClient_ReadFile_Success(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-read-file.txt"
	testData := []byte("Hello, World! This is test data.")

	// First write a file
	versionID, err := client.WriteFile(ctx, bucketName, testFileName, testData)
	if err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if versionID == "" {
		t.Error("WriteFile() should return a version ID")
	}

	// Read the file (latest version)
	readData, err := client.ReadFile(ctx, bucketName, testFileName, "")
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	// Verify data matches
	if string(readData) != string(testData) {
		t.Errorf("Expected data %s, got %s", string(testData), string(readData))
	}
}

func TestClient_ReadFile_EmptyFileName(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	_, err := client.ReadFile(ctx, bucketName, "", "")
	if err == nil {
		t.Error("ReadFile() should have failed with empty file name")
	}
}

func TestClient_ReadFile_EmptyBucketName(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	_, err := client.ReadFile(ctx, "", "test.txt", "")
	if err == nil {
		t.Error("ReadFile() should have failed with empty bucket name")
	}
}

func TestClient_ReadFile_NonExistentFile(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	_, err := client.ReadFile(ctx, bucketName, "non-existent-file.txt", "")
	if err == nil {
		t.Error("ReadFile() should have failed with non-existent file")
	}
}

func TestClient_WriteFile_Success(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-write-file.txt"
	testData := []byte("This is test data for writing.")

	versionID, err := client.WriteFile(ctx, bucketName, testFileName, testData)
	if err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if versionID == "" {
		t.Error("WriteFile() should return a version ID")
	}

	// Verify by reading back (latest version)
	readData, err := client.ReadFile(ctx, bucketName, testFileName, "")
	if err != nil {
		t.Fatalf("ReadFile() failed after WriteFile: %v", err)
	}

	if string(readData) != string(testData) {
		t.Errorf("Expected data %s, got %s", string(testData), string(readData))
	}
}

func TestClient_WriteFile_EmptyFileName(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	_, err := client.WriteFile(ctx, bucketName, "", []byte("test"))
	if err == nil {
		t.Error("WriteFile() should have failed with empty file name")
	}
}

func TestClient_WriteFile_EmptyBucketName(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	_, err := client.WriteFile(ctx, "", "test.txt", []byte("test"))
	if err == nil {
		t.Error("WriteFile() should have failed with empty bucket name")
	}
}

func TestClient_WriteFile_NilData(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	_, err := client.WriteFile(ctx, bucketName, "test.txt", nil)
	if err == nil {
		t.Error("WriteFile() should have failed with nil data")
	}
}

func TestClient_WriteFile_EmptyData(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-empty-file.txt"
	versionID, err := client.WriteFile(ctx, bucketName, testFileName, []byte{})
	if err != nil {
		t.Fatalf("WriteFile() should succeed with empty data, got: %v", err)
	}
	if versionID == "" {
		t.Error("WriteFile() should return a version ID")
	}

	// Verify by reading back (latest version)
	readData, err := client.ReadFile(ctx, bucketName, testFileName, "")
	if err != nil {
		t.Fatalf("ReadFile() failed after WriteFile: %v", err)
	}

	if len(readData) != 0 {
		t.Errorf("Expected empty data, got %d bytes", len(readData))
	}
}

func TestClient_WriteFile_ReadFile_RoundTrip(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-roundtrip.txt"
	testData := []byte("Round trip test data with special chars: !@#$%^&*()")

	// Write
	versionID, err := client.WriteFile(ctx, bucketName, testFileName, testData)
	if err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if versionID == "" {
		t.Error("WriteFile() should return a version ID")
	}

	// Read (latest version)
	readData, err := client.ReadFile(ctx, bucketName, testFileName, "")
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	// Verify
	if string(readData) != string(testData) {
		t.Errorf("Round trip failed: expected %s, got %s", string(testData), string(readData))
	}
}

func TestClient_CreateBucket_Success(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	bucketName := "test-create-bucket"

	err := client.CreateBucket(ctx, bucketName)
	if err != nil {
		t.Fatalf("CreateBucket() failed: %v", err)
	}

	// Verify bucket exists
	exists, err := client.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		t.Fatalf("BucketExists() failed: %v", err)
	}
	if !exists {
		t.Error("Bucket should exist after CreateBucket()")
	}

	// Verify versioning is enabled
	versioning, err := client.minioClient.GetBucketVersioning(ctx, bucketName)
	if err != nil {
		t.Fatalf("GetBucketVersioning() failed: %v", err)
	}
	if versioning.Status != "Enabled" {
		t.Errorf("Expected versioning status 'Enabled', got '%s'", versioning.Status)
	}
}

func TestClient_CreateBucket_EmptyBucketName(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	err := client.CreateBucket(ctx, "")
	if err == nil {
		t.Error("CreateBucket() should have failed with empty bucket name")
	}
}

func TestClient_CreateBucket_AlreadyExists(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	bucketName := "test-existing-bucket"

	// Create bucket first time
	err := client.CreateBucket(ctx, bucketName)
	if err != nil {
		t.Fatalf("CreateBucket() failed on first call: %v", err)
	}

	// Create bucket second time (should not error)
	err = client.CreateBucket(ctx, bucketName)
	if err != nil {
		t.Fatalf("CreateBucket() should succeed when bucket already exists, got: %v", err)
	}

	// Verify versioning is still enabled
	versioning, err := client.minioClient.GetBucketVersioning(ctx, bucketName)
	if err != nil {
		t.Fatalf("GetBucketVersioning() failed: %v", err)
	}
	if versioning.Status != "Enabled" {
		t.Errorf("Expected versioning status 'Enabled', got '%s'", versioning.Status)
	}
}

func TestClient_CreateBucket_WithVersioning(t *testing.T) {
	client, _ := setupMockClient(t)

	ctx := context.Background()
	bucketName := "test-versioning-bucket"

	err := client.CreateBucket(ctx, bucketName)
	if err != nil {
		t.Fatalf("CreateBucket() failed: %v", err)
	}

	// Verify versioning is enabled
	versioning, err := client.minioClient.GetBucketVersioning(ctx, bucketName)
	if err != nil {
		t.Fatalf("GetBucketVersioning() failed: %v", err)
	}
	if versioning.Status != "Enabled" {
		t.Errorf("Expected versioning status 'Enabled', got '%s'", versioning.Status)
	}
}

func TestClient_WriteFile_ReadFile_WithVersionID(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-versioned-file.txt"

	// Write first version
	firstData := []byte("First version of the file")
	firstVersionID, err := client.WriteFile(ctx, bucketName, testFileName, firstData)
	if err != nil {
		t.Fatalf("WriteFile() failed for first version: %v", err)
	}
	if firstVersionID == "" {
		t.Error("WriteFile() should return a version ID for first version")
	}

	// Write second version
	secondData := []byte("Second version of the file")
	secondVersionID, err := client.WriteFile(ctx, bucketName, testFileName, secondData)
	if err != nil {
		t.Fatalf("WriteFile() failed for second version: %v", err)
	}
	if secondVersionID == "" {
		t.Error("WriteFile() should return a version ID for second version")
	}
	if firstVersionID == secondVersionID {
		t.Error("Different writes should return different version IDs")
	}

	// Read latest version (should be second version)
	latestData, err := client.ReadFile(ctx, bucketName, testFileName, "")
	if err != nil {
		t.Fatalf("ReadFile() failed for latest version: %v", err)
	}
	if string(latestData) != string(secondData) {
		t.Errorf("Expected latest version to be second version. Expected: %s, got: %s", string(secondData), string(latestData))
	}

	// Read first version by version ID
	firstVersionData, err := client.ReadFile(ctx, bucketName, testFileName, firstVersionID)
	if err != nil {
		t.Fatalf("ReadFile() failed for first version: %v", err)
	}
	if string(firstVersionData) != string(firstData) {
		t.Errorf("Expected first version data. Expected: %s, got: %s", string(firstData), string(firstVersionData))
	}

	// Read second version by version ID
	secondVersionData, err := client.ReadFile(ctx, bucketName, testFileName, secondVersionID)
	if err != nil {
		t.Fatalf("ReadFile() failed for second version: %v", err)
	}
	if string(secondVersionData) != string(secondData) {
		t.Errorf("Expected second version data. Expected: %s, got: %s", string(secondData), string(secondVersionData))
	}
}

func TestClient_ReadFile_InvalidVersionID(t *testing.T) {
	client, bucketName := setupMockClient(t)

	ctx := context.Background()
	testFileName := "test-file.txt"
	testData := []byte("Test data")

	// Write a file
	_, err := client.WriteFile(ctx, bucketName, testFileName, testData)
	if err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}

	// Try to read with invalid version ID
	_, err = client.ReadFile(ctx, bucketName, testFileName, "invalid-version-id")
	if err == nil {
		t.Error("ReadFile() should have failed with invalid version ID")
	}
}
