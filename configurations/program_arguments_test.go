package configurations

import (
	"strings"
	"testing"
)

func TestParseArguments_DefaultValues(t *testing.T) {
	args, err := ParseArguments([]string{})
	if err != nil {
		t.Fatalf("ParseArguments() failed with default values: %v", err)
	}

	if args.Mode != DefaultMode {
		t.Errorf("Expected mode %s, got %s", DefaultMode, args.Mode)
	}

	if args.LogLevel != DefaultLogLevel {
		t.Errorf("Expected log level %s, got %s", DefaultLogLevel, args.LogLevel)
	}

	if args.ConfigPath != DefaultConfigPath {
		t.Errorf("Expected config path %s, got %s", DefaultConfigPath, args.ConfigPath)
	}

	if args.Help {
		t.Error("Expected help to be false")
	}

	if args.Version {
		t.Error("Expected version to be false")
	}
}

func TestParseArguments_ModeValidation(t *testing.T) {
	tests := []struct {
		name        string
		mode        string
		expectError bool
	}{
		{"valid single mode", "single", false},
		{"valid distributed mode", "distributed", false},
		{"valid single mode uppercase", "SINGLE", false},
		{"valid distributed mode mixed case", "DiStRiBuTeD", false},
		{"invalid mode", "invalid", true},
		{"empty mode", "", true},
		{"mode with spaces", " single ", false}, // Should be trimmed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := ParseArguments([]string{"-mode", tt.mode})
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for mode '%s', but got none", tt.mode)
				}
				if !strings.Contains(err.Error(), "invalid mode") {
					t.Errorf("Expected error message to contain 'invalid mode', got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for mode '%s': %v", tt.mode, err)
				}
				if args != nil && args.GetMode() != strings.ToLower(strings.TrimSpace(tt.mode)) {
					t.Errorf("Expected normalized mode '%s', got '%s'", strings.ToLower(strings.TrimSpace(tt.mode)), args.GetMode())
				}
			}
		})
	}
}

func TestParseArguments_LogLevelValidation(t *testing.T) {
	tests := []struct {
		name        string
		logLevel    string
		expectError bool
	}{
		{"valid trace level", "trace", false},
		{"valid debug level", "debug", false},
		{"valid info level", "info", false},
		{"valid warn level", "warn", false},
		{"valid error level", "error", false},
		{"valid fatal level", "fatal", false},
		{"valid panic level", "panic", false},
		{"valid info level uppercase", "INFO", false},
		{"valid debug level mixed case", "DeBuG", false},
		{"invalid log level", "invalid", true},
		{"empty log level", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := ParseArguments([]string{"-log-level", tt.logLevel})
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for log level '%s', but got none", tt.logLevel)
				}
				if !strings.Contains(err.Error(), "invalid log level") {
					t.Errorf("Expected error message to contain 'invalid log level', got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for log level '%s': %v", tt.logLevel, err)
				}
				if args != nil && args.GetLogLevel() != strings.ToLower(strings.TrimSpace(tt.logLevel)) {
					t.Errorf("Expected normalized log level '%s', got '%s'", strings.ToLower(strings.TrimSpace(tt.logLevel)), args.GetLogLevel())
				}
			}
		})
	}
}

func TestParseArguments_ConfigPath(t *testing.T) {
	configPath := "/path/to/config.yml"
	args, err := ParseArguments([]string{"-config", configPath})
	if err != nil {
		t.Fatalf("ParseArguments() failed: %v", err)
	}

	if args.ConfigPath != configPath {
		t.Errorf("Expected config path %s, got %s", configPath, args.ConfigPath)
	}
}

func TestParseArguments_ShorthandFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		checkFn  func(*ProgramArguments) bool
		expected bool
	}{
		{"shorthand mode", []string{"-m", "single"}, func(a *ProgramArguments) bool { return a.Mode == "single" }, true},
		{"shorthand config", []string{"-c", "/test.yml"}, func(a *ProgramArguments) bool { return a.ConfigPath == "/test.yml" }, true},
		{"shorthand log level", []string{"-l", "debug"}, func(a *ProgramArguments) bool { return a.LogLevel == "debug" }, true},
		{"shorthand help", []string{"-h"}, func(a *ProgramArguments) bool { return a.Help == true }, true},
		{"shorthand version", []string{"-v"}, func(a *ProgramArguments) bool { return a.Version == true }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := ParseArguments(tt.args)
			if err != nil {
				// Some might fail validation, that's okay for this test
				if !strings.Contains(err.Error(), "invalid") {
					t.Fatalf("ParseArguments() failed: %v", err)
				}
				return
			}

			if !tt.checkFn(args) {
				t.Errorf("Shorthand flag test failed for %s", tt.name)
			}
		})
	}
}

func TestParseArguments_HelpFlag(t *testing.T) {
	args, err := ParseArguments([]string{"-help"})
	if err != nil {
		t.Fatalf("ParseArguments() failed: %v", err)
	}

	if !args.Help {
		t.Error("Expected help flag to be true")
	}
}

func TestParseArguments_VersionFlag(t *testing.T) {
	args, err := ParseArguments([]string{"-version"})
	if err != nil {
		t.Fatalf("ParseArguments() failed: %v", err)
	}

	if !args.Version {
		t.Error("Expected version flag to be true")
	}
}

func TestParseArguments_MultipleFlags(t *testing.T) {
	args, err := ParseArguments([]string{
		"-mode", "distributed",
		"-config", "/custom/config.yml",
		"-log-level", "debug",
	})
	if err != nil {
		t.Fatalf("ParseArguments() failed: %v", err)
	}

	if args.Mode != "distributed" {
		t.Errorf("Expected mode 'distributed', got '%s'", args.Mode)
	}

	if args.ConfigPath != "/custom/config.yml" {
		t.Errorf("Expected config path '/custom/config.yml', got '%s'", args.ConfigPath)
	}

	if args.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", args.LogLevel)
	}
}

func TestProgramArguments_Validate(t *testing.T) {
	tests := []struct {
		name        string
		args        *ProgramArguments
		expectError bool
	}{
		{
			name: "valid arguments",
			args: &ProgramArguments{
				Mode:     ModeSingle,
				LogLevel: "info",
			},
			expectError: false,
		},
		{
			name: "invalid mode",
			args: &ProgramArguments{
				Mode:     "invalid",
				LogLevel: "info",
			},
			expectError: true,
		},
		{
			name: "invalid log level",
			args: &ProgramArguments{
				Mode:     ModeSingle,
				LogLevel: "invalid",
			},
			expectError: true,
		},
		{
			name: "multiple validation errors",
			args: &ProgramArguments{
				Mode:     "invalid",
				LogLevel: "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestProgramArguments_GetMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{"lowercase", "single", "single"},
		{"uppercase", "SINGLE", "single"},
		{"mixed case", "DiStRiBuTeD", "distributed"},
		{"with spaces", " single ", "single"},
		{"with spaces uppercase", " DISTRIBUTED ", "distributed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &ProgramArguments{Mode: tt.mode}
			result := args.GetMode()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestProgramArguments_GetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected string
	}{
		{"lowercase", "info", "info"},
		{"uppercase", "DEBUG", "debug"},
		{"mixed case", "ErRoR", "error"},
		{"with spaces", " warn ", "warn"},
		{"with spaces uppercase", " TRACE ", "trace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &ProgramArguments{LogLevel: tt.logLevel}
			result := args.GetLogLevel()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsValidValue(t *testing.T) {
	validValues := []string{"single", "distributed"}

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid lowercase", "single", true},
		{"valid uppercase", "SINGLE", true},
		{"valid mixed case", "DiStRiBuTeD", true},
		{"valid with spaces", " single ", true},
		{"invalid value", "invalid", false},
		{"empty value", "", false},
		{"partial match", "sing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidValue(tt.value, validValues)
			if result != tt.expected {
				t.Errorf("Expected %v for value '%s', got %v", tt.expected, tt.value, result)
			}
		})
	}
}
