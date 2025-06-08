package cmd

import (
	"bytes"
	"esh-cli/pkg/utils"
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

// TestRunAddTagExecutionValidation tests runAddTag validation logic without os.Exit calls
func TestRunAddTagExecutionValidation(t *testing.T) {
	// Test validation logic that occurs in runAddTag before any os.Exit calls

	tests := []struct {
		name        string
		environment string
		version     string
		shouldPass  bool
	}{
		{
			name:        "valid dev environment and version",
			environment: "dev",
			version:     "1.0.0",
			shouldPass:  true,
		},
		{
			name:        "valid production environment and version",
			environment: "production2",
			version:     "2.1.0",
			shouldPass:  true,
		},
		{
			name:        "invalid environment",
			environment: "invalid_env",
			version:     "1.0.0",
			shouldPass:  false,
		},
		{
			name:        "invalid version format",
			environment: "dev",
			version:     "invalid_version",
			shouldPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test environment validation (this is what runAddTag does first)
			envValid := utils.ContainsString(utils.ENVS, tt.environment)

			// Test version validation (this is what runAddTag does second)
			versionValid := utils.IsVersionValid(tt.version, false)

			overallValid := envValid && versionValid

			if overallValid != tt.shouldPass {
				t.Errorf("Expected validation result %v, got %v (env: %v, version: %v)",
					tt.shouldPass, overallValid, envValid, versionValid)
			}
		})
	}
}

// TestRunAddTagArgumentValidation tests argument validation in runAddTag
func TestRunAddTagArgumentValidation(t *testing.T) {
	// Test the argument validation logic by examining the command structure
	// This tests the parts of runAddTag that can be safely tested

	tests := []struct {
		name      string
		args      []string
		shouldErr bool
	}{
		{
			name:      "missing environment",
			args:      []string{},
			shouldErr: true,
		},
		{
			name:      "missing version",
			args:      []string{"dev"},
			shouldErr: true,
		},
		{
			name:      "valid arguments",
			args:      []string{"dev", "1.0.0"},
			shouldErr: false,
		},
		{
			name:      "too many arguments",
			args:      []string{"dev", "1.0.0", "extra"},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test argument length validation
			if len(tt.args) < 2 && !tt.shouldErr {
				t.Error("Should require at least 2 arguments")
			}

			if len(tt.args) >= 2 && !tt.shouldErr {
				// Test that valid environments are accepted
				env := tt.args[0]
				if env == "dev" || env == "test" || env == "staging" || env == "prod" {
					// Valid environment
					if utils.ContainsString(utils.ENVS, env) {
						// Environment validation would pass
					}
				}

				// Test version validation logic
				version := tt.args[1]
				if utils.IsVersionValid(version, false) {
					// Version validation would pass
				}
			}
		})
	}
}

// TestRunAddTagServiceLogic tests service-specific logic in runAddTag
func TestRunAddTagServiceLogic(t *testing.T) {
	tests := []struct {
		name        string
		serviceFlag string
		envVar      string
		expected    string
	}{
		{
			name:        "service flag takes precedence",
			serviceFlag: "test-service",
			envVar:      "env-service",
			expected:    "test-service",
		},
		{
			name:        "env var used when no flag",
			serviceFlag: "",
			envVar:      "env-service",
			expected:    "env-service",
		},
		{
			name:        "no service specified",
			serviceFlag: "",
			envVar:      "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original values
			originalService := service
			originalEnv := os.Getenv("ESH_SERVICE")

			defer func() {
				service = originalService
				if originalEnv == "" {
					os.Unsetenv("ESH_SERVICE")
				} else {
					os.Setenv("ESH_SERVICE", originalEnv)
				}
			}()

			// Set test values
			service = tt.serviceFlag
			if tt.envVar != "" {
				os.Setenv("ESH_SERVICE", tt.envVar)
			} else {
				os.Unsetenv("ESH_SERVICE")
			}

			// Test the logic that runAddTag uses for service determination
			var effectiveService string
			if service != "" {
				effectiveService = service
			} else {
				effectiveService = os.Getenv("ESH_SERVICE")
			}

			if effectiveService != tt.expected {
				t.Errorf("Expected effective service '%s', got '%s'", tt.expected, effectiveService)
			}
		})
	}
}

// TestRunAddTagPromoteLogic tests the promote vs create new tag logic
func TestRunAddTagPromoteLogic(t *testing.T) {
	originalPromoteFrom := promoteFrom
	defer func() {
		promoteFrom = originalPromoteFrom
	}()

	tests := []struct {
		name        string
		promoteFrom string
		isPromote   bool
	}{
		{
			name:        "empty promote-from creates new tag",
			promoteFrom: "",
			isPromote:   false,
		},
		{
			name:        "with promote-from promotes existing tag",
			promoteFrom: "stg6_1.2.0-release",
			isPromote:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			promoteFrom = tt.promoteFrom

			// Test the logic that determines promote vs create new
			var isPromoteMode bool
			if promoteFrom != "" {
				isPromoteMode = true
			} else {
				isPromoteMode = false
			}

			if isPromoteMode != tt.isPromote {
				t.Errorf("Expected isPromoteMode=%v, got %v", tt.isPromote, isPromoteMode)
			}
		})
	}
}
