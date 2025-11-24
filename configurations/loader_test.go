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
	yamlContent := `shardCount: 5
nats:
  url: nats://localhost:4222`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("SHARD_COUNT")
	os.Unsetenv("NATS_URL")
	defer os.Unsetenv("NATS_URL")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.ShardCount != 5 {
		t.Errorf("Expected ShardCount to be 5, got %d", cfg.ShardCount)
	}

	// Verify NATS URL is loaded from YAML
	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Errorf("Expected NATS URL to be 'nats://localhost:4222', got %s", cfg.NATS.URL)
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

	// Should have zero value for ShardCount
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

	// Check that NATS_URL maps to nats.url
	expectedNATSKey := "nats.url"
	if koanfKey, ok := envMap["NATS_URL"]; !ok {
		t.Error("NATS_URL should be in the env map")
	} else if koanfKey != expectedNATSKey {
		t.Errorf("Expected NATS_URL to map to %s, got %s", expectedNATSKey, koanfKey)
	}
}

func TestLoad_NATSDefaults(t *testing.T) {
	// Test that both NATS URL and SubjectPrefix default when not provided
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("NATS_URL")
	os.Unsetenv("NATS_SUBJECT_PREFIX")
	defer os.Unsetenv("NATS_URL")
	defer os.Unsetenv("NATS_SUBJECT_PREFIX")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// NATS URL should default to "nats://localhost:4222"
	if cfg.NATS.URL != DefaultNATSURL {
		t.Errorf("Expected NATS URL to be %s, got %s", DefaultNATSURL, cfg.NATS.URL)
	}

	// SubjectPrefix should default to "nimbus"
	if cfg.NATS.SubjectPrefix != DefaultNATSSubjectPrefix {
		t.Errorf("Expected SubjectPrefix to be %s, got %s", DefaultNATSSubjectPrefix, cfg.NATS.SubjectPrefix)
	}
}

func TestLoad_NATSSubjectPrefixFromYAML(t *testing.T) {
	// Test that SubjectPrefix can be set via YAML file
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5
nats:
  subjectPrefix: custom-prefix`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("NATS_SUBJECT_PREFIX")
	defer os.Unsetenv("NATS_SUBJECT_PREFIX")

	// Load config
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// SubjectPrefix should be from YAML
	if cfg.NATS.SubjectPrefix != "custom-prefix" {
		t.Errorf("Expected SubjectPrefix to be 'custom-prefix', got %s", cfg.NATS.SubjectPrefix)
	}

	// NATS URL should default
	if cfg.NATS.URL != DefaultNATSURL {
		t.Errorf("Expected NATS URL to default to %s, got %s", DefaultNATSURL, cfg.NATS.URL)
	}
}

func TestLoad_NATSUrlDefault(t *testing.T) {
	// Test that NATS URL defaults to localhost:4222 when not provided
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("NATS_URL")
	os.Unsetenv("NATS_SUBJECT_PREFIX")
	defer os.Unsetenv("NATS_URL")
	defer os.Unsetenv("NATS_SUBJECT_PREFIX")

	// Load config should succeed with default URL
	cfg, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() should succeed with default NATS URL, got error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() should return a config")
	}
	if cfg.NATS.URL != DefaultNATSURL {
		t.Errorf("Expected NATS URL to default to %s, got %s", DefaultNATSURL, cfg.NATS.URL)
	}
}

func TestLoad_NATSSubjectPrefixInvalidCharacters(t *testing.T) {
	// Test that SubjectPrefix with invalid characters fails validation
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "test_config.yml")
	yamlContent := `shardCount: 5
nats:
  subjectPrefix: "invalid prefix with spaces"`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("NATS_SUBJECT_PREFIX")
	defer os.Unsetenv("NATS_SUBJECT_PREFIX")

	// Load config should fail because SubjectPrefix contains spaces
	cfg, err := Load(yamlFile)
	if err == nil {
		t.Error("Load() should have failed with invalid SubjectPrefix, but didn't")
	}
	if cfg != nil {
		t.Error("Load() should return nil config on error")
	}
}

func TestLoad_NATSSubjectPrefixValidCharacters(t *testing.T) {
	// Test that SubjectPrefix with valid characters works
	testCases := []struct {
		name    string
		prefix  string
		wantErr bool
	}{
		{"alphanumeric", "nimbus123", false},
		{"with dots", "nimbus.db", false},
		{"with underscores", "nimbus_db", false},
		{"with dashes", "nimbus-db", false},
		{"with colons", "nimbus:db", false},
		{"mixed valid", "nimbus.db_v1-test:prod", false},
		{"with spaces", "nimbus db", true},
		{"with special chars", "nimbus@db", true},
		{"with slashes", "nimbus/db", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			yamlFile := filepath.Join(tmpDir, "test_config.yml")
			yamlContent := `shardCount: 5
nats:
  subjectPrefix: "` + tc.prefix + `"`
			if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
				t.Fatalf("Failed to create test YAML file: %v", err)
			}

			// Clear any existing env vars
			os.Unsetenv("NATS_SUBJECT_PREFIX")
			defer os.Unsetenv("NATS_SUBJECT_PREFIX")

			// Load config
			cfg, err := Load(yamlFile)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Load() should have failed with prefix '%s', but didn't", tc.prefix)
				}
				if cfg != nil {
					t.Error("Load() should return nil config on error")
				}
			} else {
				if err != nil {
					t.Errorf("Load() should have succeeded with prefix '%s', but failed: %v", tc.prefix, err)
				}
				if cfg != nil && cfg.NATS.SubjectPrefix != tc.prefix {
					t.Errorf("Expected SubjectPrefix to be '%s', got '%s'", tc.prefix, cfg.NATS.SubjectPrefix)
				}
			}
		})
	}
}
