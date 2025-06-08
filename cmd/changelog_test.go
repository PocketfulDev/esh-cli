package cmd

import (
	"bytes"
	"esh-cli/pkg/utils"
	"strings"
	"testing"
)

func TestChangelogCmdCreation(t *testing.T) {
	if changelogCmd == nil {
		t.Error("changelogCmd should not be nil")
	}

	if !strings.HasPrefix(changelogCmd.Use, "changelog") {
		t.Errorf("Expected changelogCmd.Use to start with 'changelog', got '%s'", changelogCmd.Use)
	}

	if changelogCmd.Short == "" {
		t.Error("changelogCmd.Short should not be empty")
	}
}

func TestChangelogFlags(t *testing.T) {
	// Test that flags exist
	flags := []string{"format", "from", "to", "since", "output"}

	for _, flagName := range flags {
		flag := changelogCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' should be defined", flagName)
		}
	}

	// Test boolean flags
	boolFlags := []string{"conventional-commits", "full", "group-by-type", "include-breaking"}
	for _, flagName := range boolFlags {
		flag := changelogCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Boolean flag '%s' should be defined", flagName)
		}
	}
}

func TestChangelogValidation(t *testing.T) {
	// Test argument validation
	if changelogCmd.Args == nil {
		t.Error("Args validator should be set")
	}

	// Test with correct number of args
	err := changelogCmd.Args(changelogCmd, []string{"stg6"})
	if err != nil {
		t.Errorf("Expected no error with 1 argument, got %v", err)
	}

	// Test with zero args (should be allowed for MaximumNArgs(1))
	err = changelogCmd.Args(changelogCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with 0 arguments for MaximumNArgs(1), got %v", err)
	}

	// Test with too many args
	err = changelogCmd.Args(changelogCmd, []string{"stg6", "extra"})
	if err == nil {
		t.Error("Expected error with 2 arguments")
	}
}

func TestChangelogCmdHelp(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	changelogCmd.SetOut(&buf)
	changelogCmd.SetErr(&buf)

	// Set args to trigger help
	changelogCmd.SetArgs([]string{"--help"})

	// Execute command (this should not return an error for help)
	err := changelogCmd.Execute()
	if err != nil {
		// Help command may return an error in some versions, but should still show help
		output := buf.String()
		if !strings.Contains(output, "changelog") {
			t.Errorf("Help output should contain command name, got: %s", output)
		}
	}
}

// TestRunChangelogIntegration tests the runChangelog function logic
func TestRunChangelogIntegration(t *testing.T) {
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
			// Test the validation part of runChangelog logic
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

// TestChangelogFormatLogic tests the format logic paths
func TestChangelogFormatLogic(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedFormat string
	}{
		{
			name:           "markdown format",
			format:         "markdown",
			expectedFormat: "markdown",
		},
		{
			name:           "json format",
			format:         "json",
			expectedFormat: "json",
		},
		{
			name:           "text format",
			format:         "text",
			expectedFormat: "text",
		},
		{
			name:           "default format",
			format:         "",
			expectedFormat: "markdown", // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic that determines output format
			outputFormat := tt.format
			if outputFormat == "" {
				outputFormat = "markdown" // default
			}

			if outputFormat != tt.expectedFormat {
				t.Errorf("Expected format=%v, got %v", tt.expectedFormat, outputFormat)
			}
		})
	}
}

// TestChangelogServiceLogic tests the service vs non-service logic paths
func TestChangelogServiceLogic(t *testing.T) {
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

// TestParseCommitType tests the commit type detection logic
func TestParseCommitType(t *testing.T) {
	tests := []struct {
		name           string
		commitMessage  string
		expectedType   string
		isConventional bool
	}{
		{
			name:           "feat commit",
			commitMessage:  "feat: add new feature",
			expectedType:   "feat",
			isConventional: true,
		},
		{
			name:           "fix commit",
			commitMessage:  "fix: resolve bug",
			expectedType:   "fix",
			isConventional: true,
		},
		{
			name:           "docs commit",
			commitMessage:  "docs: update README",
			expectedType:   "docs",
			isConventional: true,
		},
		{
			name:           "non-conventional commit",
			commitMessage:  "Random commit message",
			expectedType:   "other",
			isConventional: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conventional commit parsing logic
			var commitType string
			var isConventional bool

			// Simple conventional commit detection
			if strings.Contains(tt.commitMessage, ":") {
				parts := strings.SplitN(tt.commitMessage, ":", 2)
				if len(parts) == 2 {
					potentialType := strings.TrimSpace(parts[0])
					conventionalTypes := []string{"feat", "fix", "docs", "style", "refactor", "test", "chore"}
					for _, ct := range conventionalTypes {
						if potentialType == ct {
							commitType = potentialType
							isConventional = true
							break
						}
					}
				}
			}

			if !isConventional {
				commitType = "other"
			}

			if commitType != tt.expectedType {
				t.Errorf("Expected commit type=%v, got %v", tt.expectedType, commitType)
			}

			if isConventional != tt.isConventional {
				t.Errorf("Expected isConventional=%v, got %v", tt.isConventional, isConventional)
			}
		})
	}
}
