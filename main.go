package main

import (
	"NimbusDb/configurations"
	"os"

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
		log.Info().Msgf("%s version: %s", configurations.AppName, "dev")
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
	log.Info().Msgf("Configuration loaded: %+v", cfg)
	if args.Mode == configurations.ModeSingle {
		log.Info().Msgf("%s is running in single mode", configurations.AppName)
	} else if args.Mode == configurations.ModeDistributed {
		log.Info().Msgf("%s is running in distributed mode", configurations.AppName)
	}
	log.Info().Msgf("%s is running and accepting requests", configurations.AppName)
}
