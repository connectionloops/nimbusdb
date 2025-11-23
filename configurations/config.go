package configurations

type Config struct {
	ShardCount uint16 `koanf:"shardCount" env:"SHARD_COUNT"`
}
