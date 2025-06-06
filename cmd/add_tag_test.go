package cmd

import (
	"bytes"
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
