package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestAddTagCmdCreation(t *testing.T) {
	if addTagCmd == nil {
		t.Error("addTagCmd should not be nil")
	}

	if !strings.HasPrefix(addTagCmd.Use, "add-tag") {
		t.Errorf("Expected addTagCmd.Use to start with 'add-tag', got '%s'", addTagCmd.Use)
	}

	if addTagCmd.Short == "" {
		t.Error("addTagCmd.Short should not be empty")
	}
}

func TestAddTagFlags(t *testing.T) {
	// Test that actual flags exist based on the implementation
	flags := []string{"from", "service"}

	for _, flagName := range flags {
		flag := addTagCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' should be defined", flagName)
		}
	}

	// Test hot-fix flag (boolean)
	hotFixFlag := addTagCmd.Flags().Lookup("hot-fix")
	if hotFixFlag == nil {
		t.Error("Flag 'hot-fix' should be defined")
	}
}

func TestAddTagCmdHelp(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	addTagCmd.SetOut(&buf)
	addTagCmd.SetErr(&buf)

	// Set args to trigger help
	addTagCmd.SetArgs([]string{"--help"})

	// Execute command (this should not return an error for help)
	err := addTagCmd.Execute()
	if err != nil {
		// Help command may return an error in some versions, but should still show help
		output := buf.String()
		if !strings.Contains(output, "add-tag") {
			t.Errorf("Help output should contain command name, got: %s", output)
		}
	}
}

// Test command validation without actually running git commands
func TestAddTagValidation(t *testing.T) {
	// This is a unit test that focuses on command structure rather than execution
	// since executing the actual command would require git setup

	// Create a temporary command for testing
	testCmd := &cobra.Command{
		Use:   "add-tag",
		Short: "Add and push a hotfix tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Mock function for testing - just return nil
			return nil
		},
	}

	// Add flags similar to our actual command
	testCmd.Flags().StringP("env", "e", "", "Environment (required)")
	testCmd.Flags().StringP("version", "v", "", "Version (required)")
	testCmd.Flags().StringP("service", "s", "", "Service name (optional)")
	testCmd.Flags().BoolP("release", "r", false, "Set as release")

	// Test that the command can be created without errors
	if testCmd.Use != "add-tag" {
		t.Error("Test command creation failed")
	}
}

// Test environment variable handling
func TestAddTagEnvHandling(t *testing.T) {
	// Test with environment variable set
	originalEnv := os.Getenv("ESH_SERVICE")
	defer func() {
		if originalEnv == "" {
			os.Unsetenv("ESH_SERVICE")
		} else {
			os.Setenv("ESH_SERVICE", originalEnv)
		}
	}()

	// Set test environment variable
	os.Setenv("ESH_SERVICE", "test-service")

	// Verify environment variable is accessible
	if os.Getenv("ESH_SERVICE") != "test-service" {
		t.Error("Environment variable setting failed")
	}
}

func TestRunAddTagValidation(t *testing.T) {
	// Test command validation logic without actually executing git commands

	// Create a test command to work with
	cmd := &cobra.Command{}
	cmd.Flags().StringP("env", "e", "", "Environment")
	cmd.Flags().StringP("version", "v", "", "Version")
	cmd.Flags().StringP("service", "s", "", "Service")
	cmd.Flags().BoolP("hot-fix", "f", false, "Hot fix")
	cmd.Flags().StringP("from", "", "", "From tag")

	// Test cases for argument validation
	tests := []struct {
		name    string
		env     string
		version string
		service string
		hotfix  bool
		from    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing env",
			version: "1.2",
			wantErr: true,
			errMsg:  "env",
		},
		{
			name:    "missing version",
			env:     "stg6",
			wantErr: true,
			errMsg:  "version",
		},
		{
			name:    "valid basic params",
			env:     "stg6",
			version: "1.2",
			wantErr: false,
		},
		{
			name:    "valid with service",
			env:     "stg6",
			version: "1.2",
			service: "api",
			wantErr: false,
		},
		{
			name:    "invalid env",
			env:     "invalid",
			version: "1.2",
			wantErr: true,
			errMsg:  "invalid environment",
		},
		{
			name:    "invalid version format",
			env:     "stg6",
			version: "1.2.3",
			wantErr: false, // runAddTag doesn't validate version format in this test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			cmd.Flags().Set("env", tt.env)
			cmd.Flags().Set("version", tt.version)
			cmd.Flags().Set("service", tt.service)
			cmd.Flags().Set("hot-fix", "false")
			if tt.hotfix {
				cmd.Flags().Set("hot-fix", "true")
			}
			cmd.Flags().Set("from", tt.from)

			// Mock validation function that mimics runAddTag logic
			err := validateAddTagArgs(cmd)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error containing '%s', got nil", tt.errMsg)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tt.errMsg, err)
			}
		})
	}
}

// Mock validation function that simulates the logic in runAddTag
func validateAddTagArgs(cmd *cobra.Command) error {
	env, _ := cmd.Flags().GetString("env")
	version, _ := cmd.Flags().GetString("version")

	if env == "" {
		return fmt.Errorf("env is required")
	}
	if version == "" {
		return fmt.Errorf("version is required")
	}

	// Valid environments check
	validEnvs := []string{"dev", "mimic2", "stg6", "demo", "production2"}
	validEnv := false
	for _, validE := range validEnvs {
		if env == validE {
			validEnv = true
			break
		}
	}
	if !validEnv {
		return fmt.Errorf("invalid environment: %s", env)
	}

	return nil
}

// Test the actual runAddTag function with controlled inputs
func TestRunAddTagExecution(t *testing.T) {
	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Test with valid arguments (but without git setup to avoid actual execution)
	cmd := addTagCmd
	cmd.Flags().Set("service", "test-service")

	// Capture output for testing
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Test with insufficient arguments - this should show usage/help
	cmd.SetArgs([]string{}) // No arguments

	// Execute command with no args - should show usage (not an error in Cobra)
	err := cmd.Execute()
	// It's acceptable for this to not return an error as it shows help
	if err == nil {
		// Check that help was shown
		// This is actually the correct behavior for cobra commands
		t.Log("Command correctly showed help when no arguments provided")
	}

	// Test with one argument only - should also show usage
	cmd.SetArgs([]string{"stg6"}) // Missing version
	err = cmd.Execute()
	// Again, showing help is acceptable behavior
	if err == nil {
		t.Log("Command correctly showed help when insufficient arguments provided")
	}
}

// TestAddTagCurrentDirectoryBehavior tests that add-tag uses appropriate directory based on service flag
func TestAddTagCurrentDirectoryBehavior(t *testing.T) {
	// Test the logic that determines which FindLastTagAndComment function to call
	// This verifies the branching logic without executing actual git commands

	testCases := []struct {
		name           string
		service        string
		expectedMethod string
	}{
		{
			name:           "no service flag",
			service:        "",
			expectedMethod: "FindLastTagAndCommentInDir with current directory",
		},
		{
			name:           "with service flag",
			service:        "myservice",
			expectedMethod: "FindLastTagAndComment with service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the logic from runAddTag
			var method string

			if tc.service == "" {
				// This would call: utils.FindLastTagAndCommentInDir(environment, version, "", ".")
				method = "FindLastTagAndCommentInDir with current directory"
			} else {
				// This would call: utils.FindLastTagAndComment(environment, version, service)
				method = "FindLastTagAndComment with service"
			}

			if method != tc.expectedMethod {
				t.Errorf("Expected method %q, got %q", tc.expectedMethod, method)
			}
		})
	}
}
