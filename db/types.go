package db

// ShardsResponse represents the response for shard count queries.
type ShardsResponse struct {
	ShardCount uint16 `json:"shardCount"`
}
