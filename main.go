package main

import (
	"NimbusDb/configurations"
	"NimbusDb/health"
	"NimbusDb/version"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
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

	// Setup logger with the specified log level
	configurations.SetupLoggerWithLevel(args.GetLogLevel())

	log.Info().Msgf("Starting %s in %s mode...", configurations.AppName, args.GetMode())

	// Load configuration
	configPath := args.ConfigPath
	if configPath == "" {
		configPath = configurations.DefaultConfigPath
	}

	cfg, err := configurations.Load(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	log.Info().
		Uint16("shardCount", cfg.ShardCount).
		Int("healthPort", cfg.HealthPort).
		Str("blobEndpoint", cfg.Blob.Endpoint).
		Bool("blobUseSSL", cfg.Blob.UseSSL).
		Msg("Configuration loaded")
	if args.Mode == configurations.ModeSingle {
		log.Info().Msgf("%s is running in single mode", configurations.AppName)
	} else if args.Mode == configurations.ModeDistributed {
		log.Info().Msgf("%s is running in distributed mode", configurations.AppName)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start health check server (port is set in config, defaults to 8080)
	health.StartHealthServer(ctx, cfg.HealthPort)

	// Mark application as ready after initialization
	health.SetReady(true)
	log.Info().Msgf("%s is running and accepting requests", configurations.AppName)

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Shutting down...")
	health.SetReady(false)
	cancel()
}
