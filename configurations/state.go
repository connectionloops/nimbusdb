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
// Performance Note: This method creates a defensive copy on every call to ensure
// immutability and prevent external modifications to the internal state. This is
// safe for current usage patterns (initialization-time calls in StartShardHandlers),
// but could impact performance if called frequently in hot paths (e.g., per-request
// or in tight loops). For performance-critical scenarios where immutability is
// guaranteed by the caller, consider using GetShardIDsUnsafe() instead.
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

// GetShardIDsUnsafe returns a direct reference to the shard IDs slice without copying.
//
// WARNING: This method is provided for performance-critical paths where allocation
// overhead is significant. The caller MUST guarantee that the returned slice will
// not be modified. Modifying the returned slice will corrupt the internal state.
// In most cases, you should use GetShardIDs() instead, which returns a safe copy.
//
// Use this method only when:
//   - The code is in a verified hot path (profiling shows GetShardIDs copy as bottleneck)
//   - The caller only performs read operations (iteration, length checks, index access)
//   - The returned slice lifetime is limited to the current function scope
//
// return:
//   - []uint16: Direct reference to the internal shard IDs slice (DO NOT MODIFY)
func (s *State) GetShardIDsUnsafe() []uint16 {
	if s == nil {
		return nil
	}
	return s.shardIDs
}
