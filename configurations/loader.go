package configurations

import (
	"os"
	"reflect"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// buildEnvToKoanfMap builds a map from env var names to koanf keys using struct tags
func buildEnvToKoanfMap() map[string]string {
	envMap := make(map[string]string)
	cfgType := reflect.TypeOf(Config{})

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		envTag := field.Tag.Get("env")
		koanfTag := field.Tag.Get("koanf")

		if envTag != "" && koanfTag != "" {
			envMap[envTag] = koanfTag
		}
	}

	return envMap
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
	return cfg, nil
}
