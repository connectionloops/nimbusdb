package db

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const (
	// PointWrite represents a write operation.
	// See devdocs/api.md (Operation Types) for details.
	PointWrite = 0
	// PointRead represents a read operation.
	// See devdocs/api.md (Operation Types) for details.
	PointRead = 1
	// CollectionWrite represents a write operation.
	// See devdocs/api.md (Operation Types) for details.
	CollectionWrite = 2
	// CollectionRead represents a read operation.
	// See devdocs/api.md (Operation Types) for details.
	CollectionRead = 3
)

// ShardHandlerInfo holds subscription and channel information for a shard handler.
type ShardHandlerInfo struct {
	Subscription *nats.Subscription
	Channel      chan *nats.Msg
}

// StartShardHandlers initializes and starts all NATS shard operation handlers.
// It subscribes to shard operation subjects for the shards this node owns.
// Currently subscribes to all shards (0 to shardCount-1) as a placeholder
// until shard ownership is implemented via raft metadata cluster.
// Panics if InitializeGlobals has not been called first.
//
// return:
//   - []*ShardHandlerInfo: All handler info created for shard handlers (includes channels for cleanup)
func StartShardHandlers() []*ShardHandlerInfo {
	if globalConfig == nil || globalNATSConn == nil || globalBlobClient == nil {
		log.Fatal().Msg("InitializeGlobals must be called before StartShardHandlers")
	}

	state := GetGlobalState()
	if state == nil {
		log.Fatal().Msg("Global state must be initialized before StartShardHandlers")
	}

	shardIDs := state.GetShardIDs()
	if len(shardIDs) == 0 {
		log.Warn().Msg("No shards assigned to this node")
		return nil
	}

	handlers := make([]*ShardHandlerInfo, 0, len(shardIDs))

	// Subscribe to each shard operation subject
	for _, shardID := range shardIDs {
		subject := fmt.Sprintf("%s.shards.%d.op", globalConfig.NATS.SubjectPrefix, shardID)
		ch := make(chan *nats.Msg, globalConfig.Db.ChannelBufferSize)
		sub, err := globalNATSConn.ChanSubscribe(subject, ch)
		if err != nil {
			log.Fatal().Err(err).Uint16("shardID", shardID).Msg("Failed to subscribe to shard operation subject")
		}

		// Start handler goroutine for this shard's channel to handle the messages
		go handleShardOperation(shardID, ch)

		handlers = append(handlers, &ShardHandlerInfo{
			Subscription: sub,
			Channel:      ch,
		})

		log.Info().Uint16("shardID", shardID).Str("subject", subject).Msg("Subscribed to shard operation subject")
	}
	return handlers
}

// handleShardOperation handles requests for shard operations (write/read).
// It processes the operation based on the type header and responds accordingly.
// params:
//   - shardID: The shard ID for this operation
//   - ch: The channel to receive the messages from
func handleShardOperation(shardID uint16, ch chan *nats.Msg) {
	for msg := range ch {
		// Extract headers
		headers, err := ExtractShardOperationHeaders(msg)
		if err != nil {
			RespondWithNatsError(msg, ErrorCodeBadRequest, err.Error())
			continue
		}

		// Route to appropriate handler based on operation type
		switch headers.OperationType {
		case PointWrite:
			handleWriteOperation(msg, shardID, headers.FileName, headers.BucketName, headers.Overwrite)
		case PointRead:
			handleReadOperation(msg, shardID, headers.FileName, headers.BucketName)
		case CollectionWrite:
			RespondWithNatsError(msg, ErrorCodeBadRequest, "collection write operation not yet implemented")
		case CollectionRead:
			RespondWithNatsError(msg, ErrorCodeBadRequest, "collection read operation not yet implemented")
		default:
			RespondWithNatsError(msg, ErrorCodeBadRequest, fmt.Sprintf("unknown operation type: %d", headers.OperationType))
		}

	}
}

// handleWriteOperation handles write requests for shard operations.
// It writes the message data directly to blob storage without parsing.
// If overwrite is false and the file already exists, it returns an error.
// params:
//   - msg: The NATS message which contains pure byte[] data to be written to blob storage
//   - shardID: The shard ID for this operation
//   - fileName: The file path where the data should be stored
//   - bucketName: The bucket name where the data should be stored
//   - overwrite: If false, returns an error if the file already exists
func handleWriteOperation(msg *nats.Msg, shardID uint16, fileName string, bucketName string, overwrite bool) {
	// todo: metrics for write latency and count
	// Create context with timeout for blob operation using config value
	ctx, cancel := context.WithTimeout(context.Background(), globalConfig.Blob.BlobOperationTimeout)
	defer cancel()

	// Check if file exists when overwrite is false
	if !overwrite {
		exists, err := globalBlobClient.FileExists(ctx, bucketName, fileName)
		if err != nil {
			RespondWithNatsError(msg, ErrorCodeInternalServerError, fmt.Sprintf("failed to check if file exists: %v", err))
			return
		}
		if exists {
			RespondWithNatsError(msg, ErrorCodeBadRequest, fmt.Sprintf("file already exists: %s", fileName))
			return
		}
	}

	// Write data directly to blob without parsing (as per API spec)
	_, err := globalBlobClient.WriteFile(ctx, bucketName, fileName, msg.Data)
	if err != nil {
		log.Error().Err(err).Str("fileName", fileName).Str("bucketName", bucketName).Uint16("shardID", shardID).Msg("Failed to write file to blob storage")
		RespondWithNatsError(msg, ErrorCodeInternalServerError, fmt.Sprintf("failed to write file: %v", err))
		return
	}

	// Respond with success
	RespondWithNatsSuccess(msg)
}

// handleReadOperation handles read requests for shard operations.
// It reads the file data directly from blob storage and returns it as byte[].
// The data is returned directly without parsing, as per API specification.
// params:
//   - msg: The NATS message to respond to
//   - shardID: The shard ID for this operation
//   - fileName: The file path to read from
//   - bucketName: The bucket name where the file is stored
func handleReadOperation(msg *nats.Msg, shardID uint16, fileName string, bucketName string) {
	// todo: metrics for read latency and count
	// Create context with timeout for blob operation using config value
	ctx, cancel := context.WithTimeout(context.Background(), globalConfig.Blob.BlobOperationTimeout)
	defer cancel()

	// Read data directly from blob without parsing (as per API spec)
	data, err := globalBlobClient.ReadFile(ctx, bucketName, fileName, "")
	if err != nil {
		log.Error().Err(err).Str("fileName", fileName).Str("bucketName", bucketName).Uint16("shardID", shardID).Msg("Failed to read file from blob storage")
		RespondWithNatsError(msg, ErrorCodeInternalServerError, fmt.Sprintf("failed to read file: %v", err))
		return
	}

	// Respond with raw byte[] data directly (as per API spec: shard owner never parses data)
	msg.Respond(data)
}
