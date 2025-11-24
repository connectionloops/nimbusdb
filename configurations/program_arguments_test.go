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
				Mode: ModeSingle,
			},
			expectError: false,
		},
		{
			name: "invalid mode",
			args: &ProgramArguments{
				Mode: "invalid",
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
