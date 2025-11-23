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

	AppName = "NimbusDb"
)
