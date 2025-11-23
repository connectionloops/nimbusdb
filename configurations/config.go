package configurations

type Config struct {
	ShardCount uint16 `koanf:"shardCount" env:"SHARD_COUNT"`
	HealthPort int    `koanf:"healthPort" env:"HEALTH_PORT"`
	Blob       BlobConfig
}

// BlobConfig holds the configuration for MinIO blob storage.
type BlobConfig struct {
	Endpoint                          string `koanf:"blob.endpoint" env:"BLOB_ENDPOINT"`
	AccessKeyID                       string `koanf:"blob.accessKeyID" env:"BLOB_ACCESS_KEY_ID"`
	SecretAccessKey                   string `koanf:"blob.secretAccessKey" env:"BLOB_SECRET_ACCESS_KEY"`
	UseSSL                            bool   `koanf:"blob.useSSL" env:"BLOB_USE_SSL"`
	DeleteMarkerCleanupDelayDays      int    `koanf:"blob.deleteMarkerCleanupDelayDays" env:"BLOB_DELETE_MARKER_CLEANUP_DELAY_DAYS"`            // in days, default 1
	NonCurrentVersionCleanupDelayDays int    `koanf:"blob.nonCurrentVersionCleanupDelayDays" env:"BLOB_NON_CURRENT_VERSION_CLEANUP_DELAY_DAYS"` // in days, default 1
}
