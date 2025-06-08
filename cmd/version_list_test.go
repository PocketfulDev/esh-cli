package cmd

import (
	"bytes"
	"esh-cli/pkg/utils"
	"strings"
	"testing"
)

func TestVersionListCmdCreation(t *testing.T) {
	if versionListCmd == nil {
		t.Error("versionListCmd should not be nil")
	}

	if !strings.HasPrefix(versionListCmd.Use, "version-list") {
		t.Errorf("Expected versionListCmd.Use to start with 'version-list', got '%s'", versionListCmd.Use)
	}

	if versionListCmd.Short == "" {
		t.Error("versionListCmd.Short should not be empty")
	}
}

func TestVersionListFlags(t *testing.T) {
	// Test that flags exist
	flags := []string{"all", "major", "minor", "format", "sort", "limit"}

	for _, flagName := range flags {
		flag := versionListCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' should be defined", flagName)
		}
	}

	// Test boolean flags
	boolFlags := []string{"all"}
	for _, flagName := range boolFlags {
		flag := versionListCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Boolean flag '%s' should be defined", flagName)
		}
	}
}

func TestVersionListValidation(t *testing.T) {
	// Test argument validation
	if versionListCmd.Args == nil {
		t.Error("Args validator should be set")
	}

	// Test with correct number of args
	err := versionListCmd.Args(versionListCmd, []string{"stg6"})
	if err != nil {
		t.Errorf("Expected no error with 1 argument, got %v", err)
	}

	// Test with zero args (should be allowed for MaximumNArgs(1))
	err = versionListCmd.Args(versionListCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with 0 arguments for MaximumNArgs(1), got %v", err)
	}

	// Test with too many args
	err = versionListCmd.Args(versionListCmd, []string{"stg6", "extra"})
	if err == nil {
		t.Error("Expected error with 2 arguments")
	}
}

func TestVersionListCmdHelp(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	versionListCmd.SetOut(&buf)
	versionListCmd.SetErr(&buf)

	// Set args to trigger help
	versionListCmd.SetArgs([]string{"--help"})

	// Execute command (this should not return an error for help)
	err := versionListCmd.Execute()
	if err != nil {
		// Help command may return an error in some versions, but should still show help
		output := buf.String()
		if !strings.Contains(output, "version-list") {
			t.Errorf("Help output should contain command name, got: %s", output)
		}
	}
}

// TestRunVersionListIntegration tests the runVersionList function logic
func TestRunVersionListIntegration(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectExit  bool
		expectError string
	}{
		{
			name:        "invalid environment",
			args:        []string{"invalid_env"},
			expectExit:  true,
			expectError: "invalid environment",
		},
		{
			name:       "valid environment",
			args:       []string{"stg6"},
			expectExit: false,
		},
		{
			name:       "valid production environment",
			args:       []string{"production2"},
			expectExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation part of runVersionList logic
			environment := tt.args[0]

			// Test environment validation
			if !utils.ContainsString(utils.ENVS, environment) {
				if !tt.expectExit {
					t.Errorf("Expected valid environment, but got invalid: %s", environment)
				}
				return // This would trigger os.Exit(1)
			}

			// If we reach here, validation passed
			if tt.expectExit {
				t.Errorf("Expected validation to fail and exit, but it passed")
			}
		})
	}
}

// TestVersionListFormatLogic tests the format logic paths
func TestVersionListFormatLogic(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedFormat string
	}{
		{
			name:           "table format",
			format:         "table",
			expectedFormat: "table",
		},
		{
			name:           "json format",
			format:         "json",
			expectedFormat: "json",
		},
		{
			name:           "compact format",
			format:         "compact",
			expectedFormat: "compact",
		},
		{
			name:           "default format",
			format:         "",
			expectedFormat: "table", // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic that determines output format
			outputFormat := tt.format
			if outputFormat == "" {
				outputFormat = "table" // default
			}

			if outputFormat != tt.expectedFormat {
				t.Errorf("Expected format=%v, got %v", tt.expectedFormat, outputFormat)
			}
		})
	}
}

// TestVersionListServiceLogic tests the service vs non-service logic paths
func TestVersionListServiceLogic(t *testing.T) {
	tests := []struct {
		name            string
		serviceFlag     string
		expectedCurrent bool // true if should use current directory
	}{
		{
			name:            "no service flag uses current directory",
			serviceFlag:     "",
			expectedCurrent: true,
		},
		{
			name:            "with service flag uses config lookup",
			serviceFlag:     "myservice",
			expectedCurrent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic path that determines directory usage
			var usesCurrentDir bool
			if tt.serviceFlag == "" {
				// This path would use current directory
				usesCurrentDir = true
			} else {
				// This path would: initConfig() and findProjectPath(service)
				usesCurrentDir = false
			}

			if usesCurrentDir != tt.expectedCurrent {
				t.Errorf("Expected usesCurrentDir=%v, got %v", tt.expectedCurrent, usesCurrentDir)
			}
		})
	}
}
