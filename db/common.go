package db

import (
	"NimbusDb/blob"
	"NimbusDb/configurations"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

const (
	// ErrorCodeBadRequest represents a client error (400)
	ErrorCodeBadRequest = 400
	// ErrorCodeInternalServerError represents a server error (500)
	ErrorCodeInternalServerError = 500

	SuccessCode = 200
)

var (
	// globalConfig holds the configuration for system handlers.
	// It is set once during initialization and never modified.
	globalConfig *configurations.Config
	// globalNATSConn holds the NATS connection for system handlers.
	// It is set once during initialization and never modified.
	globalNATSConn *nats.Conn
	// globalBlobClient holds the blob client for system handlers.
	// It is set once during initialization and never modified.
	globalBlobClient *blob.Client
	// globalShutdownCtx holds the shutdown context for graceful shutdown.
	// It is used to signal when the application is shutting down to prevent
	// new long-running operations from starting.
	globalShutdownCtx context.Context
	// initOnce ensures InitializeGlobals can only be called once.
	initOnce sync.Once
)

var (
	// globalState holds the runtime state of this node.
	// Uses atomic.Value for lock-free reads and thread-safe writes.
	globalState atomic.Value // *configurations.State
)

// ShardOperationHeaders contains the extracted headers from a shard operation request.
type ShardOperationHeaders struct {
	OperationType int
	FileName      string
	BucketName    string
	Overwrite     bool
}

type DbResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// InitializeGlobals sets the global configuration, NATS connection, blob client, and shutdown context for use by handlers.
// This should be called once during application startup before starting handlers.
// Subsequent calls to this function will be ignored (idempotent).
//
// params:
//   - cfg: The application configuration containing NATS and shard settings
//   - nc: The NATS connection to use for subscriptions
//   - blobClient: The blob client to use for storage operations
//   - shutdownCtx: The context that will be cancelled during graceful shutdown
func InitializeGlobals(cfg *configurations.Config, nc *nats.Conn, blobClient *blob.Client, shutdownCtx context.Context) {
	initOnce.Do(func() {
		globalConfig = cfg
		globalNATSConn = nc
		globalBlobClient = blobClient
		globalShutdownCtx = shutdownCtx
	})
}

// InitializeSingleModeState initializes the global state for single mode.
// It creates a state with all shard IDs from 0 to ShardCount-1.
// This function is thread-safe.
func InitializeSingleModeState() {
	shardIds := []uint16{}
	for shardID := uint16(0); shardID < globalConfig.ShardCount; shardID++ {
		shardIds = append(shardIds, shardID)
	}
	SetGlobalState(configurations.NewState(shardIds))
}

// GetGlobalState returns the current global state.
// This function is thread-safe for concurrent reads (lock-free).
// Multiple goroutines can call this simultaneously without blocking.
// The returned State object is effectively immutable (shardIDs field is unexported).
// This is optimized for read-heavy workloads with zero allocations.
//
// return:
//   - *configurations.State: The current global state, or nil if not initialized
func GetGlobalState() *configurations.State {
	val := globalState.Load()
	if val == nil {
		return nil
	}
	return val.(*configurations.State)
}

// SetGlobalState sets the global state to the provided state.
// This function is thread-safe for concurrent writes.
//
// params:
//   - state: The new state to set
func SetGlobalState(state *configurations.State) {
	globalState.Store(state)
}

// RespondWithNatsError responds with a NATS native error using headers.
// This provides a standardized, reusable error response format across all handlers.
// Uses NATS native error headers as per NATS best practices.
// Retryability can be determined by the error code: 4xx errors are typically not retriable,
// while 5xx errors may be retriable depending on the specific error.
// params:
//   - msg: The NATS message to respond to
//   - errorCode: The error code (e.g., "400" for bad request, "500" for server error)
//   - errorDescription: The error description/message
//
// return:
//   - nil
func RespondWithNatsError(msg *nats.Msg, status int, description string) {
	resp := DbResponse{
		Error:  description,
		Status: status,
	}
	b, _ := json.Marshal(resp)
	msg.Respond(b)
}

func RespondWithNatsSuccess(msg *nats.Msg) {
	resp := DbResponse{
		Error:  "",
		Status: SuccessCode,
	}
	b, _ := json.Marshal(resp)
	msg.Respond(b)

}

// extractShardOperationHeaders extracts and validates required headers from a NATS message.
// It extracts operation type, fileName, and bucketName from the message headers.
// Optimized for performance by using direct map access and explicit base parsing.
// params:
//   - msg: The NATS message containing the operation request
//
// return:
//   - *ShardOperationHeaders: The extracted headers
//   - error: An error if any required header is missing or invalid
func ExtractShardOperationHeaders(msg *nats.Msg) (*ShardOperationHeaders, error) {
	h := msg.Header

	// --- type ---
	opStr := h.Get("type")
	if opStr == "" {
		return nil, errors.New("missing 'type' header")
	}
	op, err := strconv.Atoi(opStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'type' header: %s", opStr)
	}

	// --- fileName ---
	fn := h.Get("fileName")
	if fn == "" {
		return nil, errors.New("missing 'fileName' header")
	}

	// --- bucketName ---
	bn := h.Get("bucketName")
	if bn == "" {
		return nil, errors.New("missing 'bucketName' header")
	}

	// --- overwrite (default true) ---
	owStr := h.Get("overwrite")
	ow := true // backward-compatible default
	if owStr != "" {
		ow, err = strconv.ParseBool(owStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'overwrite' header: %s", owStr)
		}
	}

	// return the struct pointer (single heap alloc)
	return &ShardOperationHeaders{
		OperationType: op,
		FileName:      fn,
		BucketName:    bn,
		Overwrite:     ow,
	}, nil
}
