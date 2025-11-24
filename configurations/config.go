package configurations

import "time"

// Config holds the configuration for NimbusDb.
// This is typically set via environment variables or a configuration file.
// This is primary way to supply configs and secrets to NimbusDb.
// These values are common to all nodes in NimbusDb cluster.
type Config struct {
	ShardCount uint16 `koanf:"shardCount" env:"SHARD_COUNT"`
	HealthPort int    `koanf:"healthPort" env:"HEALTH_PORT"`
	// LogLevel specifies the logging level.
	// Valid values: "trace", "debug", "info", "warn", "error", "fatal", "panic"
	LogLevel string     `koanf:"logLevel" env:"LOG_LEVEL"`
	Blob     BlobConfig `koanf:"blob"`
	NATS     NATSConfig `koanf:"nats"`
	Db       DbConfig   `koanf:"db"`
}

// NATSConfig holds the configuration for NATS.
type NATSConfig struct {
	URL                    string        `koanf:"url" env:"NATS_URL"`
	Creds                  string        `koanf:"creds" env:"NATS_CREDS"`
	SubjectPrefix          string        `koanf:"subjectPrefix" env:"NATS_SUBJECT_PREFIX"`
	NatsDrainTimeout       time.Duration `koanf:"natsDrainTimeout" env:"NATS_DRAIN_TIMEOUT"`             // timeout for NATS drain operation, default 30s
	ShutdownGracePeriod    time.Duration `koanf:"shutdownGracePeriod" env:"NATS_SHUTDOWN_GRACE_PERIOD"` // grace period to wait for in-flight messages during shutdown, default 100ms
}

// BlobConfig holds the configuration for MinIO blob storage.
type BlobConfig struct {
	Endpoint                          string        `koanf:"endpoint" env:"BLOB_ENDPOINT"`
	AccessKeyID                       string        `koanf:"accessKeyID" env:"BLOB_ACCESS_KEY_ID"`
	SecretAccessKey                   string        `koanf:"secretAccessKey" env:"BLOB_SECRET_ACCESS_KEY"`
	UseSSL                            bool          `koanf:"useSSL" env:"BLOB_USE_SSL"`
	DeleteMarkerCleanupDelayDays      int           `koanf:"deleteMarkerCleanupDelayDays" env:"BLOB_DELETE_MARKER_CLEANUP_DELAY_DAYS"`            // in days, default 1
	NonCurrentVersionCleanupDelayDays int           `koanf:"nonCurrentVersionCleanupDelayDays" env:"BLOB_NON_CURRENT_VERSION_CLEANUP_DELAY_DAYS"` // in days, default 1
	BlobOperationTimeout              time.Duration `koanf:"blobOperationTimeout" env:"BLOB_OPERATION_TIMEOUT"`                                   // timeout for blob operations, default 30s
}

type DbConfig struct {
	ChannelBufferSize int `koanf:"channelBufferSize" env:"DB_CHANNEL_BUFFER_SIZE"`
}

const (
	LogLevelTrace = "trace"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"
)
