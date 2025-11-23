package configurations

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger configures the global logger with console output and colors.
// Uses the default log level (info).
func SetupLogger() {
	SetupLoggerWithLevel("info")
}

// SetupLoggerWithLevel configures the global logger with console output and colors.
// The log level is set based on the provided level string.
//
// params:
//   - level: The log level string (trace, debug, info, warn, error, fatal, panic).
//     If an invalid level is provided, defaults to info.
func SetupLoggerWithLevel(level string) {
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse log level
	logLevel := parseLogLevel(level)
	zerolog.SetGlobalLevel(logLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false, // Enable colors
	})
}

// parseLogLevel converts a log level string to a zerolog.Level.
// If the level is invalid, returns zerolog.InfoLevel as default.
//
// params:
//   - level: The log level string (case-insensitive).
//
// return:
//   - zerolog.Level: The corresponding zerolog level.
func parseLogLevel(level string) zerolog.Level {
	levelLower := strings.ToLower(strings.TrimSpace(level))
	switch levelLower {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
