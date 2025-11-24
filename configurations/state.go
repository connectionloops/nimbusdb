package configurations

// State represents the state of this node.
// State can change mid runtime.
type State struct {
	shardIDs []uint16
}

// NewState creates a new State instance with the given shard IDs.
//
// params:
//   - shardIDs: The list of shard IDs assigned to this node
//
// return:
//   - *State: A new State instance
func NewState(shardIDs []uint16) *State {
	return &State{
		shardIDs: shardIDs,
	}
}

// GetShardIDs returns a copy of the shard IDs assigned to this node.
//
// return:
//   - []uint16: A copy of the shard IDs list
func (s *State) GetShardIDs() []uint16 {
	if s == nil {
		return nil
	}

	// Return a copy to prevent external modifications
	result := make([]uint16, len(s.shardIDs))
	copy(result, s.shardIDs)
	return result
}
