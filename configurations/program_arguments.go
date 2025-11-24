package configurations

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// ProgramArguments holds all command-line arguments for the application.
// This typically control runtime behavior of the application.
// These values are specific to each node in NimbusDb cluster.
type ProgramArguments struct {
	// Mode specifies the operation mode of the application.
	// Valid values: "single", "distributed"
	Mode string

	// ConfigPath specifies the path to the configuration YAML file.
	// Defaults to ".config.yml" if not specified.
	ConfigPath string

	// Help displays the help message and exits.
	Help bool

	// Version displays the version information and exits.
	Version bool
}

const (
	// ModeSingle represents single-node operation mode
	ModeSingle = "single"

	// ModeDistributed represents distributed operation mode
	ModeDistributed = "distributed"

	// DefaultMode is the default operation mode
	DefaultMode = ModeSingle

	// DefaultConfigPath is the default configuration file path
	DefaultConfigPath = ".config.yml"
)

var (
	validModes = []string{ModeSingle, ModeDistributed}
)

// ParseArguments parses command-line arguments and returns a ProgramArguments struct.
// It validates all arguments and returns an error if any validation fails.
//
// params:
//   - args: Command-line arguments (typically os.Args[1:]). If nil, uses flag.CommandLine.
//
// return:
//   - *ProgramArguments: The parsed and validated arguments.
//   - error: An error if parsing or validation fails.
func ParseArguments(args []string) (*ProgramArguments, error) {
	parsedArgs := &ProgramArguments{}

	// Create a new flag set to avoid conflicts with other flag usage
	fs := flag.NewFlagSet("nimbusdb", flag.ContinueOnError)

	fs.StringVar(&parsedArgs.Mode, "mode", DefaultMode, fmt.Sprintf("Operation mode: %s", strings.Join(validModes, " or ")))
	fs.StringVar(&parsedArgs.Mode, "m", DefaultMode, "Shorthand for -mode")
	fs.StringVar(&parsedArgs.ConfigPath, "config", DefaultConfigPath, fmt.Sprintf("Path to configuration YAML file (default: %s)", DefaultConfigPath))
	fs.StringVar(&parsedArgs.ConfigPath, "c", DefaultConfigPath, "Shorthand for -config")
	fs.BoolVar(&parsedArgs.Help, "help", false, "Display help message")
	fs.BoolVar(&parsedArgs.Help, "h", false, "Shorthand for -help")
	fs.BoolVar(&parsedArgs.Version, "version", false, "Display version information")
	fs.BoolVar(&parsedArgs.Version, "v", false, "Shorthand for -version")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [options]\n\n", AppName)
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	// Parse arguments
	if args == nil {
		args = flag.Args()
	}
	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate parsed arguments
	if err := parsedArgs.Validate(); err != nil {
		return nil, err
	}

	return parsedArgs, nil
}

// Validate validates all fields in ProgramArguments and returns an error if any validation fails.
//
// return:
//   - error: An error if validation fails, nil otherwise.
func (pa *ProgramArguments) Validate() error {
	var validationErrors []string

	// Validate mode
	if !isValidValue(pa.Mode, validModes) {
		validationErrors = append(validationErrors, fmt.Sprintf("invalid mode '%s': must be one of %s", pa.Mode, strings.Join(validModes, ", ")))
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return nil
}

// isValidValue checks if a value is in the list of valid values (case-insensitive).
//
// params:
//   - value: The value to check.
//   - validValues: The list of valid values.
//
// return:
//   - bool: True if the value is valid, false otherwise.
func isValidValue(value string, validValues []string) bool {
	valueLower := strings.ToLower(strings.TrimSpace(value))
	for _, valid := range validValues {
		if strings.ToLower(valid) == valueLower {
			return true
		}
	}
	return false
}

// GetMode returns the normalized mode value (lowercase, trimmed).
//
// return:
//   - string: The normalized mode value.
func (pa *ProgramArguments) GetMode() string {
	return strings.ToLower(strings.TrimSpace(pa.Mode))
}
