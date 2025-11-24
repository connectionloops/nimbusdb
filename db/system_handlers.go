package db

import (
	"NimbusDb/configurations"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// StartSystemHandlers initializes and starts all NATS system handlers.
// It subscribes to system subjects using the globally configured connection.
// Panics if InitializeGlobals has not been called first.
//
// return:
//   - []*nats.Subscription: All subscriptions created for system handlers
func StartSystemHandlers() []*nats.Subscription {
	if globalConfig == nil || globalNATSConn == nil || globalBlobClient == nil {
		log.Fatal().Msg("InitializeGlobals must be called before StartSystemHandlers")
	}

	var subscriptions []*nats.Subscription

	// Subscribe to shard count requests with the queue group as specified in API docs
	sub, err := globalNATSConn.QueueSubscribe(
		globalConfig.NATS.SubjectPrefix+".config.getShardCount",
		configurations.SystemHandlersQueueGroup,
		getShardCount,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start NATS system handler")
	}
	subscriptions = append(subscriptions, sub)

	return subscriptions
}

// getShardCount handles requests for the current shard count.
// It responds with the shard count from the global configuration.
func getShardCount(msg *nats.Msg) {
	resp := ShardsResponse{
		ShardCount: globalConfig.ShardCount,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		RespondWithNatsError(msg, ErrorCodeInternalServerError, err.Error())
		return
	}
	if err := msg.Respond(b); err != nil {
		log.Error().Err(err).Msg("Failed to respond to getShardCount")
	}
}
