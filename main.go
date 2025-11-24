package main

import (
	"NimbusDb/blob"
	"NimbusDb/configurations"
	"NimbusDb/db"
	"NimbusDb/health"
	"NimbusDb/version"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

//go:embed banner.txt
var banner string

func main() {
	// Print ASCII art banner
	fmt.Print(banner)

	// Parse command-line arguments
	args, err := configurations.ParseArguments(os.Args[1:])
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse arguments")
	}

	// Handle help flag
	if args.Help {
		// Usage is already printed by flag package
		os.Exit(0)
	}

	// Handle version flag
	if args.Version {
		log.Info().Msgf("%s version: %s", configurations.AppName, version.GetVersion())
		os.Exit(0)
	}

	log.Info().Msgf("Starting %s in %s mode...", configurations.AppName, args.GetMode())

	// Load configuration
	cfg := configurations.MustLoad(args.ConfigPath)

	// Setup logger with the specified log level
	configurations.SetupLoggerWithLevel(cfg.LogLevel)

	// setup NATS client
	nc := connectNATS(cfg)

	// setup blob client
	ctx := context.Background()
	blobClient, err := blob.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create blob client")
	}

	db.InitializeGlobals(cfg, nc, blobClient)
	systemSubscriptions := db.StartSystemHandlers()

	var shardHandlers []*db.ShardHandlerInfo

	switch args.GetMode() {
	case configurations.ModeSingle:
		db.InitializeSingleModeState()
		shardHandlers = db.StartShardHandlers()
	case configurations.ModeDistributed:
		log.Fatal().Msg("Distributed mode is not supported yet")
	default:
		log.Fatal().Msgf("Invalid mode: %s", args.GetMode())
	}

	// Collect all subscriptions for graceful shutdown
	subscriptions := make([]*nats.Subscription, 0, len(systemSubscriptions)+len(shardHandlers))
	subscriptions = append(subscriptions, systemSubscriptions...)
	for _, handler := range shardHandlers {
		subscriptions = append(subscriptions, handler.Subscription)
	}

	// Create context for graceful shutdown
	shutdownCtx, cancel := context.WithCancel(context.Background())

	// Start health check server (port is set in config, defaults to 8080)
	health.StartHealthServer(shutdownCtx, cfg.HealthPort)

	// Mark application as ready after initialization
	health.SetReady(true)
	log.Info().Msgf("%s is running and accepting requests", configurations.AppName)

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Shutting down...")

	// Mark as not ready to stop accepting new requests
	health.SetReady(false)

	// Cancel context to trigger health server shutdown
	cancel()

	// Drain NATS connection and unsubscribe from all subscriptions
	drainNats(nc, subscriptions, shardHandlers, cfg)

	log.Info().Msg("Graceful shutdown complete. bye bye!")
}

// connectNATS establishes a connection to the NATS server using the provided configuration.
// Configures production-ready connection options including reconnect handling and error callbacks.
// Reconnection is handled by NATS client automatically.
// params:
//   - cfg: The application configuration containing NATS connection settings
//
// return:
//   - *nats.Conn: The established NATS connection
func connectNATS(cfg *configurations.Config) *nats.Conn {
	nc, err := nats.Connect(
		cfg.NATS.URL,
		nats.UserCredentialBytes([]byte(cfg.NATS.Creds)),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	return nc
}

// drainNats gracefully shuts down NATS by unsubscribing from all subscriptions,
// waiting for in-flight messages to complete, closing channels, and then draining the connection with a timeout.
func drainNats(nc *nats.Conn, subscriptions []*nats.Subscription, shardHandlers []*db.ShardHandlerInfo, cfg *configurations.Config) {
	// Unsubscribe from all NATS subscriptions first to stop receiving new messages
	log.Info().Msg("Unsubscribing from NATS subjects...")
	for _, sub := range subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			log.Error().Err(err).Msg("Failed to unsubscribe from NATS subject")
		}
	}

	// Wait briefly for any in-flight messages to be processed
	log.Info().Msg("Waiting for in-flight messages to complete...")
	time.Sleep(cfg.NATS.ShutdownGracePeriod)

	// Now it's safe to close all shard handler channels
	log.Info().Msg("Closing shard handler channels...")
	for _, handler := range shardHandlers {
		close(handler.Channel)
	}

	// Drain the NATS connection to allow in-flight messages to complete
	// Use a timeout to prevent hanging indefinitely
	log.Info().Msg("Draining NATS connection...")
	drainDone := make(chan error, 1)
	go func() {
		drainDone <- nc.Drain()
	}()

	select {
	case err := <-drainDone:
		if err != nil {
			log.Error().Err(err).Msg("Failed to drain NATS connection")
		} else {
			log.Info().Msg("NATS connection drained successfully")
		}
	case <-time.After(cfg.NATS.NatsDrainTimeout):
		log.Warn().Dur("timeout", cfg.NATS.NatsDrainTimeout).Msg("NATS drain operation timed out, closing connection")
		nc.Close()
	}
}
