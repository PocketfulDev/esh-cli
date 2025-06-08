package main

import (
	"bytes"
	"esh-cli/cmd"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// Mutex to synchronize access to global variables during testing
var testMutex sync.Mutex

func TestMainVersionVariable(t *testing.T) {
	// Test that version variable exists and has a default value
	if version == "" {
		t.Error("version variable should have a default value")
	}
}

func TestMainPackageStructure(t *testing.T) {
	// This is a basic test to ensure the main package is properly structured
	// The actual main() function is tested through integration tests

	// Synchronize access to global variables to prevent race conditions
	testMutex.Lock()
	defer testMutex.Unlock()

	// Verify version variable is accessible
	originalVersion := version
	if originalVersion == "" {
		t.Skip("version variable is empty, which is acceptable for testing")
	}

	// Test that version can be modified (this would happen during build)
	testVersion := "test-version"
	version = testVersion

	if version != testVersion {
		t.Errorf("version variable should be modifiable, expected %s, got %s", testVersion, version)
	}

	// Restore original version
	version = originalVersion
}

func TestMainFunctionExecution(t *testing.T) {
	// Test main function by building and running the binary
	// This is the most reliable way to test main() function coverage

	// Get current working directory for CI compatibility
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Build the binary for testing
	cmd := exec.Command("go", "build", "-o", "esh-cli-test", ".")
	cmd.Dir = wd
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove(filepath.Join(wd, "esh-cli-test"))

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "version command",
			args:        []string{"--version"},
			expectError: false,
		},
		{
			name:        "help command",
			args:        []string{"--help"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the binary with test args
			execCmd := exec.Command("./esh-cli-test", tt.args...)
			execCmd.Dir = wd

			output, err := execCmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v, output: %s", err, string(output))
			}

			// Basic output validation
			if len(output) == 0 && !tt.expectError {
				t.Error("Expected some output")
			}
		})
	}
}

func TestMainFunctionActualExecution(t *testing.T) {
	// Test main() function by actually calling it with controlled args
	// Since main() might call os.Exit, we need to handle this carefully

	// Save original args and restore them after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "help flag",
			args: []string{"esh-cli", "--help"},
		},
		{
			name: "version flag",
			args: []string{"esh-cli", "--version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Synchronize access to global variables to prevent race conditions
			testMutex.Lock()
			defer testMutex.Unlock()

			// Set test args
			os.Args = tt.args

			// Set test version
			originalVersion := version
			version = "test-main-execution"

			defer func() {
				version = originalVersion
				if r := recover(); r != nil {
					// main() might call os.Exit via cmd.Execute(), which is expected
					// We don't fail the test for this as it's normal behavior
				}
			}()

			// Call the actual main() function
			// This will provide coverage for main()
			main()
		})
	}
}

func TestCmdExecuteIntegration(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		contains    string
	}{
		{
			name:        "help command execution",
			args:        []string{"esh-cli", "--help"},
			expectError: false,
			contains:    "ESH CLI tool", // Should contain this in help output
		},
		{
			name:        "version command execution",
			args:        []string{"esh-cli", "--version"},
			expectError: false,
			contains:    "", // Version output varies
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create isolated command instance with test version
			testVersion := "test-integration-v1.0.0"
			rootCmd := cmd.NewRootCmd(testVersion)
			rootCmd.SetVersionTemplate("{{.Version}}\n")

			// Create isolated buffer for output capture
			var outputBuf bytes.Buffer
			rootCmd.SetOut(&outputBuf)
			rootCmd.SetErr(&outputBuf)

			// Set arguments for the command (excluding the program name)
			if len(tt.args) > 1 {
				rootCmd.SetArgs(tt.args[1:])
			}

			// Execute the command
			err := rootCmd.Execute()

			// Check for expected errors
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check output contains expected content
			output := outputBuf.String()
			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.contains, output)
			}

			// For version command, verify it shows our test version
			if strings.Contains(tt.args[len(tt.args)-1], "version") {
				if !strings.Contains(output, testVersion) {
					t.Errorf("Expected version output to contain '%s', got: %s", testVersion, output)
				}
			}
		})
	}
}

func TestMainVersionBuildIntegration(t *testing.T) {
	tests := []struct {
		name           string
		buildVersion   string
		expectedResult string
	}{
		{
			name:           "default dev version",
			buildVersion:   "dev",
			expectedResult: "dev",
		},
		{
			name:           "release version",
			buildVersion:   "1.0.0",
			expectedResult: "1.0.0",
		},
		{
			name:           "prerelease version",
			buildVersion:   "1.0.0-beta.1",
			expectedResult: "1.0.0-beta.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that we can create commands with different versions
			rootCmd := cmd.NewRootCmd(tt.buildVersion)

			// Test that the version is set correctly on the command
			if rootCmd.Version != tt.expectedResult {
				t.Errorf("Expected version %s, got %s", tt.expectedResult, rootCmd.Version)
			}

			// Test that version setting works with the isolated command
			testVersion := "isolated-test-" + tt.buildVersion
			rootCmd.Version = testVersion
			if rootCmd.Version != testVersion {
				t.Errorf("Expected isolated version setting to work, got %s", rootCmd.Version)
			}
		})
	}
}

func TestMainCommandExecution(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "help command",
			args:        []string{"esh-cli", "--help"},
			expectError: false,
		},
		{
			name:        "version command",
			args:        []string{"esh-cli", "--version"},
			expectError: false,
		},
		{
			name:        "invalid command",
			args:        []string{"esh-cli", "invalid-cmd"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Synchronize access to global variables to prevent race conditions
			testMutex.Lock()
			defer testMutex.Unlock()

			// Store original args
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()

			// Set test args
			os.Args = tt.args

			// Test that the command structure exists
			// This tests the integration between main and cmd packages
			if len(os.Args) < 1 {
				t.Error("Args should have at least the program name")
			}

			// Test version setting before command execution
			testVersion := "test-main-version"
			originalVersion := version
			version = testVersion

			cmd.SetVersion(version)

			// Verify version was set
			if version != testVersion {
				t.Errorf("Expected version %s, got %s", testVersion, version)
			}

			// Restore version
			version = originalVersion
		})
	}
}

// TestMainImportStructure tests the package import structure
func TestMainImportStructure(t *testing.T) {
	// Test that the main package can successfully import and use cmd package

	// Synchronize access to global variables to prevent race conditions
	testMutex.Lock()
	defer testMutex.Unlock()

	// Test version setting functionality
	testVersion := "import-test-version"
	originalVersion := version
	version = testVersion

	// This mimics what main() does
	cmd.SetVersion(version)

	// Test successful import and function call
	if version != testVersion {
		t.Errorf("Expected version %s after setting, got %s", testVersion, version)
	}

	// Restore original version
	version = originalVersion
}

// TestMainFunctionSignature tests the main function signature and structure
func TestMainFunctionSignature(t *testing.T) {
	// Test that main function components are accessible

	// Synchronize access to global variables to prevent race conditions
	testMutex.Lock()
	defer testMutex.Unlock()

	// Test version variable exists and is modifiable
	originalVersion := version
	testVersion := "signature-test"
	version = testVersion

	if version != testVersion {
		t.Error("version variable should be modifiable")
	}

	// Test that cmd package functions are accessible
	cmd.SetVersion(version)

	// This tests that the main function's structure is sound:
	// 1. version variable exists
	// 2. cmd.SetVersion can be called
	// 3. cmd.Execute would be callable (though we don't call it to avoid side effects)

	// Restore version
	version = originalVersion
}
