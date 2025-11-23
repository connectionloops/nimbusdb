package configurations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FromYAML(t *testing.T) {
	// Create a temporary YAML file
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env var
	os.Unsetenv("SHARD_COUNT")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.ShardCount != 5 {
		t.Errorf("Expected ShardCount to be 5, got %d", cfg.ShardCount)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	// Create a temporary directory (no YAML file)
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.yml")

	// Set environment variable
	os.Setenv("SHARD_COUNT", "10")
	defer os.Unsetenv("SHARD_COUNT")

	// Load config
	cfg, err := Load(nonExistentFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.ShardCount != 10 {
		t.Errorf("Expected ShardCount to be 10, got %d", cfg.ShardCount)
	}
}

func TestLoad_EnvOverridesYAML(t *testing.T) {
	// Create a temporary YAML file
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set environment variable (should override YAML)
	os.Setenv("SHARD_COUNT", "20")
	defer os.Unsetenv("SHARD_COUNT")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Environment variable should override YAML value
	if cfg.ShardCount != 20 {
		t.Errorf("Expected ShardCount to be 20 (from env), got %d", cfg.ShardCount)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	// Clear any existing env var
	os.Unsetenv("SHARD_COUNT")

	// Load with empty path (should not error, but config will have zero values)
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() with empty path should not error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should have zero value
	if cfg.ShardCount != 0 {
		t.Errorf("Expected ShardCount to be 0 (zero value), got %d", cfg.ShardCount)
	}
}

func TestLoad_NonExistentYAMLFile(t *testing.T) {
	// Clear any existing env var
	os.Unsetenv("SHARD_COUNT")

	// Load with non-existent file path (should not error)
	cfg, err := Load("/nonexistent/path/config.yml")
	if err != nil {
		t.Fatalf("Load() with non-existent file should not error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should have zero value
	if cfg.ShardCount != 0 {
		t.Errorf("Expected ShardCount to be 0 (zero value), got %d", cfg.ShardCount)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create a temporary YAML file with invalid content
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "invalid_config.yml")
	invalidYAML := `shardCount: [invalid`
	if err := os.WriteFile(yamlFile, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env var
	os.Unsetenv("SHARD_COUNT")

	// Load should fail with invalid YAML
	cfg, err := Load(yamlFile)
	if err == nil {
		t.Error("Load() with invalid YAML should have failed, but didn't")
	}
	if cfg != nil {
		t.Error("Load() should return nil config on error")
	}
}

func TestLoad_EnvVarMapping(t *testing.T) {
	// Test that the env var name SHARD_COUNT correctly maps to koanf key shardCount
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set environment variable using the struct tag name
	os.Setenv("SHARD_COUNT", "42")
	defer os.Unsetenv("SHARD_COUNT")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Should use env var value, not YAML value
	if cfg.ShardCount != 42 {
		t.Errorf("Expected ShardCount to be 42 (from SHARD_COUNT env var), got %d", cfg.ShardCount)
	}
}

func TestLoad_ComplexScenario(t *testing.T) {
	// Test a more complex scenario: YAML exists, env var exists, then env var is removed
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 7`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// First load with env var set
	os.Setenv("SHARD_COUNT", "15")
	cfg1, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg1.ShardCount != 15 {
		t.Errorf("Expected ShardCount to be 15, got %d", cfg1.ShardCount)
	}

	// Remove env var and load again
	os.Unsetenv("SHARD_COUNT")
	cfg2, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg2.ShardCount != 7 {
		t.Errorf("Expected ShardCount to be 7 (from YAML), got %d", cfg2.ShardCount)
	}
}

func TestBuildEnvToKoanfMap(t *testing.T) {
	// Test that the mapping function correctly builds the map
	envMap := buildEnvToKoanfMap()

	// Check that SHARD_COUNT maps to shardCount
	expectedKey := "shardCount"
	if koanfKey, ok := envMap["SHARD_COUNT"]; !ok {
		t.Error("SHARD_COUNT should be in the env map")
	} else if koanfKey != expectedKey {
		t.Errorf("Expected SHARD_COUNT to map to %s, got %s", expectedKey, koanfKey)
	}
}
