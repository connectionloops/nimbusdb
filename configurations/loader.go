package configurations

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

// buildEnvToKoanfMap builds a map from env var names to koanf keys using struct tags.
// It recursively processes nested structs to handle embedded configurations.
func buildEnvToKoanfMap() map[string]string {
	envMap := make(map[string]string)
	cfgType := reflect.TypeOf(Config{})

	buildEnvToKoanfMapRecursive(cfgType, "", envMap)
	return envMap
}

// buildEnvToKoanfMapRecursive recursively processes struct fields to build the env to koanf map.
func buildEnvToKoanfMapRecursive(typ reflect.Type, prefix string, envMap map[string]string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		envTag := field.Tag.Get("env")
		koanfTag := field.Tag.Get("koanf")

		if envTag != "" && koanfTag != "" {
			envMap[envTag] = koanfTag
		}

		// Handle nested structs (but skip if it's already been processed via tags)
		if field.Type.Kind() == reflect.Struct && envTag == "" {
			// Recursively process nested struct
			buildEnvToKoanfMapRecursive(field.Type, prefix, envMap)
		}
	}
}

// Load loads the configuration from the given path.
// If the path is empty, it will load the configuration from the environment variables.
// If the path is not empty, it will load the configuration from the YAML file.
// The configuration is loaded from the YAML file first, and then the environment variables are used to override the configuration.
// The environment variables are used to override the configuration by prefixing the environment variable name with the prefix "NIMBUS_DB_".
//
// params:
//   - path: The path to the YAML configuration file. If empty, only environment variables will be used.
//
// return:
//   - *Config: The loaded configuration struct.
//   - error: An error if the configuration could not be loaded.
func Load(path string) (*Config, error) {
	cfg := &Config{}
	k := koanf.New(".")

	// 1. Load base YAML file if it exists (don't error if it doesn't)
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
				return nil, err
			}
		}
	}

	// 2. Override with env vars using struct tags
	envMap := buildEnvToKoanfMap()
	if err := k.Load(env.Provider("", ".", func(s string) string {
		// Look up the env var name in our map to get the koanf key
		if koanfKey, ok := envMap[s]; ok {
			return koanfKey
		}
		// If not found, return as-is (fallback)
		return strings.ToLower(s)
	}), nil); err != nil {
		return nil, err
	}

	// 3. Unmarshal final configurations
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	// 4. Set defaults for fields that weren't set
	if cfg.HealthPort == 0 {
		cfg.HealthPort = DefaultHealthPort
	}
	if cfg.Blob.DeleteMarkerCleanupDelayDays == 0 {
		cfg.Blob.DeleteMarkerCleanupDelayDays = DefaultDeleteMarkerCleanupDelayDays
	}
	if cfg.Blob.NonCurrentVersionCleanupDelayDays == 0 {
		cfg.Blob.NonCurrentVersionCleanupDelayDays = DefaultNonCurrentVersionCleanupDelayDays
	}

	// 5. Validate configuration
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	// 6. Log configuration
	log.Info().Msgf(`Configuration loaded:
		shardCount: %d
		healthPort: %d
		blobEndpoint: %s
		blobUseSSL: %t
		blobDeleteMarkerCleanupDelayDays: %d
		blobNonCurrentVersionCleanupDelayDays: %d`,
		cfg.ShardCount,
		cfg.HealthPort,
		cfg.Blob.Endpoint,
		cfg.Blob.UseSSL,
		cfg.Blob.DeleteMarkerCleanupDelayDays,
		cfg.Blob.NonCurrentVersionCleanupDelayDays,
	)

	return cfg, nil
}

// validateConfig validates the configuration values.
func validateConfig(cfg *Config) error {
	// Validate health port range (1-65535)
	if cfg.HealthPort < 1 || cfg.HealthPort > 65535 {
		return fmt.Errorf("health port must be between 1 and 65535, got %d", cfg.HealthPort)
	}

	// Validate lifecycle cleanup delay days (1-365 days, 1 year max)
	const maxLifecycleDays = 365
	if cfg.Blob.DeleteMarkerCleanupDelayDays < 1 || cfg.Blob.DeleteMarkerCleanupDelayDays > maxLifecycleDays {
		return fmt.Errorf("delete marker cleanup delay days must be between 1 and %d, got %d", maxLifecycleDays, cfg.Blob.DeleteMarkerCleanupDelayDays)
	}
	if cfg.Blob.NonCurrentVersionCleanupDelayDays < 1 || cfg.Blob.NonCurrentVersionCleanupDelayDays > maxLifecycleDays {
		return fmt.Errorf("non-current version cleanup delay days must be between 1 and %d, got %d", maxLifecycleDays, cfg.Blob.NonCurrentVersionCleanupDelayDays)
	}

	return nil
}
