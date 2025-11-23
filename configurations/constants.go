package configurations

const (
	// DefaultShardCount is the default number of shards
	DefaultShardCount uint16 = 16

	// MaxShardCount is the maximum allowed number of shards
	MaxShardCount uint16 = 256

	// MinShardCount is the minimum allowed number of shards
	MinShardCount uint16 = 1

	// DefaultHealthPort is the default port for the health check server
	DefaultHealthPort int = 8080

	// DefaultDeleteMarkerCleanupDelayDays is the default delay in days before delete markers are cleaned up
	DefaultDeleteMarkerCleanupDelayDays int = 1

	// DefaultNonCurrentVersionCleanupDelayDays is the default delay in days before non-current versions are cleaned up
	DefaultNonCurrentVersionCleanupDelayDays int = 1

	AppName = "NimbusDb"
)
