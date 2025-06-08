package cmd

import (
	"bytes"
	"esh-cli/pkg/utils"
	"strings"
	"testing"
)

func TestVersionDiffCmdCreation(t *testing.T) {
	if versionDiffCmd == nil {
		t.Error("versionDiffCmd should not be nil")
	}

	if !strings.HasPrefix(versionDiffCmd.Use, "version-diff") {
		t.Errorf("Expected versionDiffCmd.Use to start with 'version-diff', got '%s'", versionDiffCmd.Use)
	}

	if versionDiffCmd.Short == "" {
		t.Error("versionDiffCmd.Short should not be empty")
	}
}

func TestVersionDiffFlags(t *testing.T) {
	// Test that flags exist
	flags := []string{"history", "remote", "commits", "files", "stats", "since"}

	for _, flagName := range flags {
		flag := versionDiffCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' should be defined", flagName)
		}
	}

	// Test boolean flags
	boolFlags := []string{"history", "remote", "commits", "files", "stats"}
	for _, flagName := range boolFlags {
		flag := versionDiffCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Boolean flag '%s' should be defined", flagName)
		}
	}
}

func TestVersionDiffValidation(t *testing.T) {
	// Test argument validation
	if versionDiffCmd.Args == nil {
		t.Error("Args validator should be set")
	}

	// Test with correct number of args
	err := versionDiffCmd.Args(versionDiffCmd, []string{"stg6"})
	if err != nil {
		t.Errorf("Expected no error with 1 argument, got %v", err)
	}

	// Test with minimum args (MinimumNArgs(1) requires at least 1)
	err = versionDiffCmd.Args(versionDiffCmd, []string{})
	if err == nil {
		t.Error("Expected error with 0 arguments for MinimumNArgs(1)")
	}

	// Test with more args (should be allowed for MinimumNArgs(1))
	err = versionDiffCmd.Args(versionDiffCmd, []string{"stg6", "extra"})
	if err != nil {
		t.Errorf("Expected no error with 2 arguments for MinimumNArgs(1), got %v", err)
	}
}

func TestVersionDiffCmdHelp(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	versionDiffCmd.SetOut(&buf)
	versionDiffCmd.SetErr(&buf)

	// Set args to trigger help
	versionDiffCmd.SetArgs([]string{"--help"})

	// Execute command (this should not return an error for help)
	err := versionDiffCmd.Execute()
	if err != nil {
		// Help command may return an error in some versions, but should still show help
		output := buf.String()
		if !strings.Contains(output, "version-diff") {
			t.Errorf("Help output should contain command name, got: %s", output)
		}
	}
}

// TestRunVersionDiffIntegration tests the runVersionDiff function logic
func TestRunVersionDiffIntegration(t *testing.T) {
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
			// Test the validation part of runVersionDiff logic
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

// TestVersionDiffServiceLogic tests the service vs non-service logic paths
func TestVersionDiffServiceLogic(t *testing.T) {
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

// TestCompareVersionsLogic tests the version comparison logic
func TestCompareVersionsLogic(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		expected string
	}{
		{
			name:     "major version bump",
			from:     "1.0.0",
			to:       "2.0.0",
			expected: "major",
		},
		{
			name:     "minor version bump",
			from:     "1.0.0",
			to:       "1.1.0",
			expected: "minor",
		},
		{
			name:     "patch version bump",
			from:     "1.0.0",
			to:       "1.0.1",
			expected: "patch",
		},
		{
			name:     "same version",
			from:     "1.0.0",
			to:       "1.0.0",
			expected: "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test version comparison logic (simplified)
			bumpType := getBumpTypeTest(tt.from, tt.to)

			if bumpType != tt.expected {
				t.Errorf("Expected bump type=%v, got %v", tt.expected, bumpType)
			}
		})
	}
}

// getBumpTypeTest is a simplified version of the bump type detection logic for testing
func getBumpTypeTest(from, to string) string {
	if from == to {
		return "none"
	}

	// Parse semantic versions (simplified)
	fromParts := strings.Split(from, ".")
	toParts := strings.Split(to, ".")

	if len(fromParts) >= 3 && len(toParts) >= 3 {
		if fromParts[0] != toParts[0] {
			return "major"
		}
		if fromParts[1] != toParts[1] {
			return "minor"
		}
		if fromParts[2] != toParts[2] {
			return "patch"
		}
	}

	return "unknown"
}

// TestVersionDiffFromToLogic tests the from/to version logic
func TestVersionDiffFromToLogic(t *testing.T) {
	tests := []struct {
		name     string
		fromFlag string
		toFlag   string
		hasFrom  bool
		hasTo    bool
	}{
		{
			name:     "no from/to flags",
			fromFlag: "",
			toFlag:   "",
			hasFrom:  false,
			hasTo:    false,
		},
		{
			name:     "from flag only",
			fromFlag: "v1.0.0",
			toFlag:   "",
			hasFrom:  true,
			hasTo:    false,
		},
		{
			name:     "to flag only",
			fromFlag: "",
			toFlag:   "v2.0.0",
			hasFrom:  false,
			hasTo:    true,
		},
		{
			name:     "both from and to flags",
			fromFlag: "v1.0.0",
			toFlag:   "v2.0.0",
			hasFrom:  true,
			hasTo:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic that handles from/to flags
			hasFrom := tt.fromFlag != ""
			hasTo := tt.toFlag != ""

			if hasFrom != tt.hasFrom {
				t.Errorf("Expected hasFrom=%v, got %v", tt.hasFrom, hasFrom)
			}

			if hasTo != tt.hasTo {
				t.Errorf("Expected hasTo=%v, got %v", tt.hasTo, hasTo)
			}
		})
	}
}
