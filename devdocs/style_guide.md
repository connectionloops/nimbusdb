# Style Guide

This document outlines the coding standards and conventions for the NimbusDb project.

## Naming Conventions

- **Folder names** (packages): Use lowercase, short names without underscores. Multi-word names should be concatenated (e.g., `httputil`, `filepath`). Prefer singular forms.
- **File names**: Use `snake_case`
- **Variables**: Use `camelCase`
- **Methods**:
  - `PascalCase` for exported methods
  - `camelCase` for unexported methods
- **Structs**:
  - `PascalCase` for exported structs
  - `camelCase` for unexported structs

> Note: go automatically decides if a method (or struct) is exported or not based on first letter being upper case or lower case.

## Comment Conventions

All exported functions, structs, and fields should have documentation comments following this format:

### Function Comments

```go
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
    // ...
}
```

**Format guidelines:**

- Start with a sentence describing what the function does (beginning with the function name)
- Add additional sentences explaining behavior, edge cases, or important details
- Use a blank line (`//`) to separate the description from the params/return sections
- Document all parameters in a `params:` section with bullet points
- Document all return values in a `return:` section with bullet points
- Use single-line comments (`//`) for documentation comments

### Struct Comments

```go
// Config holds the application configuration loaded from YAML files
// and environment variables. Environment variables take precedence
// over values in the YAML file.
type Config struct {
    // ShardCount specifies the number of database shards to use.
    // Must be between MinShardCount and MaxShardCount.
    ShardCount uint16 `koanf:"shardCount" env:"SHARD_COUNT"`
}
```

**Format guidelines:**

- Start with the struct name and describe what it represents
- Document fields with comments placed directly above each field
- Explain constraints, defaults, or important behavior for fields
