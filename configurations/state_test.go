package configurations

import (
	"testing"
)

// TestNewState verifies that NewState creates a State instance with the given shard IDs.
func TestNewState(t *testing.T) {
	shardIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(shardIDs)

	if state == nil {
		t.Fatal("NewState returned nil")
	}

	if state.shardIDs == nil {
		t.Fatal("State shardIDs is nil")
	}

	if len(state.shardIDs) != len(shardIDs) {
		t.Errorf("Expected %d shard IDs, got %d", len(shardIDs), len(state.shardIDs))
	}

	for i, id := range shardIDs {
		if state.shardIDs[i] != id {
			t.Errorf("Expected shard ID %d at index %d, got %d", id, i, state.shardIDs[i])
		}
	}
}

// TestNewStateEmptySlice verifies that NewState can handle an empty shard ID slice.
func TestNewStateEmptySlice(t *testing.T) {
	state := NewState([]uint16{})

	if state == nil {
		t.Fatal("NewState returned nil")
	}

	if state.shardIDs == nil {
		t.Fatal("State shardIDs is nil")
	}

	if len(state.shardIDs) != 0 {
		t.Errorf("Expected 0 shard IDs, got %d", len(state.shardIDs))
	}
}

// TestNewStateDefensiveCopy verifies that NewState creates a defensive copy of the input slice.
func TestNewStateDefensiveCopy(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	// Modify the original slice
	originalIDs[0] = 999
	originalIDs[2] = 888

	// Verify the state's internal slice was not affected
	stateIDs := state.GetShardIDs()
	if stateIDs[0] != 1 {
		t.Errorf("Expected internal state to be unaffected: expected 1, got %d", stateIDs[0])
	}
	if stateIDs[2] != 3 {
		t.Errorf("Expected internal state to be unaffected: expected 3, got %d", stateIDs[2])
	}
}

// TestGetShardIDs verifies that GetShardIDs returns a copy of the shard IDs.
func TestGetShardIDs(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	retrievedIDs := state.GetShardIDs()

	if retrievedIDs == nil {
		t.Fatal("GetShardIDs returned nil")
	}

	if len(retrievedIDs) != len(originalIDs) {
		t.Errorf("Expected %d shard IDs, got %d", len(originalIDs), len(retrievedIDs))
	}

	for i, id := range originalIDs {
		if retrievedIDs[i] != id {
			t.Errorf("Expected shard ID %d at index %d, got %d", id, i, retrievedIDs[i])
		}
	}
}

// TestGetShardIDsImmutability verifies that modifying the returned slice doesn't affect internal state.
func TestGetShardIDsImmutability(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	retrievedIDs := state.GetShardIDs()
	// Modify the returned slice
	retrievedIDs[0] = 999
	retrievedIDs[2] = 888

	// Get the shard IDs again and verify they're unchanged
	unchangedIDs := state.GetShardIDs()

	for i, id := range originalIDs {
		if unchangedIDs[i] != id {
			t.Errorf("Internal state was modified: expected %d at index %d, got %d", id, i, unchangedIDs[i])
		}
	}
}

// TestGetShardIDsNilState verifies that GetShardIDs handles nil state gracefully.
func TestGetShardIDsNilState(t *testing.T) {
	var state *State = nil
	result := state.GetShardIDs()

	if result != nil {
		t.Errorf("Expected nil for nil state, got %v", result)
	}
}

// TestGetShardIDsEmptyState verifies that GetShardIDs handles empty shard IDs.
func TestGetShardIDsEmptyState(t *testing.T) {
	state := NewState([]uint16{})
	result := state.GetShardIDs()

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestGetShardIDsUnsafe verifies that GetShardIDsUnsafe returns a direct reference.
func TestGetShardIDsUnsafe(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	retrievedIDs := state.GetShardIDsUnsafe()

	if retrievedIDs == nil {
		t.Fatal("GetShardIDsUnsafe returned nil")
	}

	if len(retrievedIDs) != len(originalIDs) {
		t.Errorf("Expected %d shard IDs, got %d", len(originalIDs), len(retrievedIDs))
	}

	for i, id := range originalIDs {
		if retrievedIDs[i] != id {
			t.Errorf("Expected shard ID %d at index %d, got %d", id, i, retrievedIDs[i])
		}
	}
}

// TestGetShardIDsUnsafeNilState verifies that GetShardIDsUnsafe handles nil state gracefully.
func TestGetShardIDsUnsafeNilState(t *testing.T) {
	var state *State = nil
	result := state.GetShardIDsUnsafe()

	if result != nil {
		t.Errorf("Expected nil for nil state, got %v", result)
	}
}

// TestGetShardIDsUnsafeEmptyState verifies that GetShardIDsUnsafe handles empty shard IDs.
func TestGetShardIDsUnsafeEmptyState(t *testing.T) {
	state := NewState([]uint16{})
	result := state.GetShardIDsUnsafe()

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestGetShardIDsUnsafeIsReference verifies that GetShardIDsUnsafe returns a direct reference
// (not a copy) by checking that modifications affect the internal state.
// NOTE: This test demonstrates the "unsafe" nature - in production code, you should NEVER
// modify the returned slice.
func TestGetShardIDsUnsafeIsReference(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	unsafeIDs := state.GetShardIDsUnsafe()
	// Store original value
	originalValue := unsafeIDs[0]

	// Modify via unsafe reference (DON'T DO THIS IN REAL CODE)
	unsafeIDs[0] = 999

	// Verify internal state was affected
	unsafeIDsAgain := state.GetShardIDsUnsafe()
	if unsafeIDsAgain[0] != 999 {
		t.Errorf("Expected internal state to be modified to 999, got %d", unsafeIDsAgain[0])
	}

	// Restore for consistency
	unsafeIDs[0] = originalValue
}

// TestGetShardIDsVsUnsafe verifies that the safe version returns a copy while unsafe returns a reference.
func TestGetShardIDsVsUnsafe(t *testing.T) {
	originalIDs := []uint16{1, 2, 3, 4, 5}
	state := NewState(originalIDs)

	safeIDs := state.GetShardIDs()
	unsafeIDs := state.GetShardIDsUnsafe()

	// Both should have the same values initially
	for i := range originalIDs {
		if safeIDs[i] != unsafeIDs[i] {
			t.Errorf("Values differ at index %d: safe=%d, unsafe=%d", i, safeIDs[i], unsafeIDs[i])
		}
	}

	// Modify safe copy - should not affect internal state
	safeIDs[0] = 888

	// Verify internal state is unchanged via unsafe reference
	unsafeIDsCheck := state.GetShardIDsUnsafe()
	if unsafeIDsCheck[0] != originalIDs[0] {
		t.Errorf("Safe copy modification affected internal state: expected %d, got %d", originalIDs[0], unsafeIDsCheck[0])
	}

	// Also verify via safe getter
	safeIDsCheck := state.GetShardIDs()
	if safeIDsCheck[0] != originalIDs[0] {
		t.Errorf("Safe copy modification affected internal state: expected %d, got %d", originalIDs[0], safeIDsCheck[0])
	}
}
