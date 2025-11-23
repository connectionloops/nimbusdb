package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the version of the application, set at build time
	Version = "dev"
	// Commit is the git commit SHA, set at build time
	Commit = "unknown"
	// BuildDate is the build date, set at build time
	BuildDate = "unknown"
	// GoVersion is the Go version used to build the application
	GoVersion = runtime.Version()
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns a detailed version string
func GetFullVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, go: %s)", Version, Commit, BuildDate, GoVersion)
}

// GetShortVersion returns a short version string
func GetShortVersion() string {
	return Version
}
