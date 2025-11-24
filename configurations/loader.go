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

		// Build the full koanf key path
		var fullKoanfKey string
		if prefix != "" {
			if koanfTag != "" {
				fullKoanfKey = prefix + "." + koanfTag
			} else {
				fullKoanfKey = prefix
			}
		} else {
			fullKoanfKey = koanfTag
		}

		if envTag != "" && fullKoanfKey != "" {
			envMap[envTag] = fullKoanfKey
		}

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct && koanfTag != "" {
			// Recursively process nested struct with the parent's koanf tag as prefix
			buildEnvToKoanfMapRecursive(field.Type, fullKoanfKey, envMap)
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
	if cfg.NATS.URL == "" {
		cfg.NATS.URL = DefaultNATSURL
	}
	if cfg.NATS.SubjectPrefix == "" {
		cfg.NATS.SubjectPrefix = DefaultNATSSubjectPrefix
	}
	if cfg.NATS.NatsDrainTimeout == 0 {
		cfg.NATS.NatsDrainTimeout = DefaultNATSDrainTimeout
	}
	if cfg.Blob.BlobOperationTimeout == 0 {
		cfg.Blob.BlobOperationTimeout = DefaultBlobOperationTimeout
	}
	if cfg.Db.ChannelBufferSize == 0 {
		cfg.Db.ChannelBufferSize = DefaultDbChannelBufferSize
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DefaultLogLevel
	}
	// 5. Validate configuration
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	// 6. Log configuration
	log.Info().Msg("Configuration loaded:")
	log.Info().Msgf("shardCount: %d", cfg.ShardCount)
	log.Info().Msgf("healthPort: %d", cfg.HealthPort)
	log.Info().Msgf("blobEndpoint: %s", cfg.Blob.Endpoint)
	log.Info().Msgf("blobUseSSL: %t", cfg.Blob.UseSSL)
	log.Info().Msgf("blobDeleteMarkerCleanupDelayDays: %d", cfg.Blob.DeleteMarkerCleanupDelayDays)
	log.Info().Msgf("blobNonCurrentVersionCleanupDelayDays: %d", cfg.Blob.NonCurrentVersionCleanupDelayDays)
	log.Info().Msgf("natsURL: %s", cfg.NATS.URL)
	log.Info().Msgf("natsSubjectPrefix: %s", cfg.NATS.SubjectPrefix)
	log.Info().Msgf("dbChannelBufferSize: %d", cfg.Db.ChannelBufferSize)
	log.Info().Msgf("logLevel: %s", cfg.LogLevel)

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

	// Validate NATS configuration
	if err := validateNATSConfig(&cfg.NATS); err != nil {
		return err
	}

	return nil
}

// MustLoad loads the configuration from the given path and calls log.Fatal() if an error occurs.
// This is a convenience function for use in main() to keep error handling out of the main flow.
//
// params:
//   - path: The path to the YAML configuration file. If empty, only environment variables will be used.
//
// return:
//   - *Config: The loaded configuration struct. The function will never return nil as it calls log.Fatal() on error.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	return cfg
}

// validateNATSConfig validates the NATS configuration values.
func validateNATSConfig(cfg *NATSConfig) error {
	// SubjectPrefix must be a valid NATS subject prefix
	// NATS subject tokens can contain alphanumeric characters, dots, underscores, dashes, and colons
	// They cannot contain spaces or be empty
	if cfg.SubjectPrefix == "" {
		return fmt.Errorf("NATS subject prefix cannot be empty")
	}

	// Check for invalid characters in subject prefix
	// Valid characters: alphanumeric, dots, underscores, dashes, colons
	for _, char := range cfg.SubjectPrefix {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' ||
			char == '_' ||
			char == '-' ||
			char == ':') {
			return fmt.Errorf("NATS subject prefix contains invalid character '%c': only alphanumeric characters, dots, underscores, dashes, and colons are allowed", char)
		}
	}

	return nil
}
