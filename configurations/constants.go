package configurations

import "time"

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

	// DefaultNATSSubjectPrefix is the default subject prefix for NATS
	DefaultNATSSubjectPrefix string = "nimbus"

	// DefaultNATSURL is the default NATS server URL
	DefaultNATSURL string = "nats://localhost:4222"

	// DefaultDbChannelBufferSize is the default channel buffer size for the database operations
	DefaultDbChannelBufferSize int = 256

	// DefaultLogLevel is the default logging level
	DefaultLogLevel string = LogLevelInfo

	// DefaultNATSDrainTimeout is the default timeout for NATS drain operation
	DefaultNATSDrainTimeout = 30 * time.Second

	// DefaultNATSShutdownGracePeriod is the default grace period to wait for in-flight messages during shutdown
	DefaultNATSShutdownGracePeriod = 100 * time.Millisecond

	// DefaultBlobOperationTimeout is the default timeout for blob operations
	DefaultBlobOperationTimeout = 30 * time.Second

	AppName = "NimbusDb"

	SystemHandlersQueueGroup = "common_config_qg"
)
