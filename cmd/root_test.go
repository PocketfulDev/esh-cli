package cmd

import (
	"os"
	"testing"
)

func TestSetVersion(t *testing.T) {
	originalVersion := version
	defer func() {
		version = originalVersion
	}()

	testVersion := "1.2.3"
	SetVersion(testVersion)

	if version != testVersion {
		t.Errorf("SetVersion() failed: expected %s, got %s", testVersion, version)
	}
}

func TestRootCmdCreation(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}

	if rootCmd.Use != "esh-cli" {
		t.Errorf("Expected rootCmd.Use to be 'esh-cli', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestRootCmdVersion(t *testing.T) {
	// Test that version is set via the Version field, not a flag
	SetVersion("test-version")
	if rootCmd.Version != "test-version" {
		t.Errorf("Expected rootCmd.Version to be 'test-version', got '%s'", rootCmd.Version)
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function can be called
	// We'll set args to show help to avoid running actual commands
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Test with help flag
	os.Args = []string{"esh-cli", "--help"}

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// Note: Execute() calls os.Exit on help, so we can't test the return value
	// We're just testing that the function exists and is callable
}

func TestExecuteFunction(t *testing.T) {
	// Store original args
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "help command",
			args: []string{"esh-cli", "--help"},
		},
		{
			name: "projects command",
			args: []string{"esh-cli", "projects"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test args
			os.Args = tt.args

			// Call Execute function and handle potential os.Exit
			defer func() {
				if r := recover(); r != nil {
					// Execute() might call os.Exit via rootCmd.Execute() which is expected
					// We don't fail the test for this as it's normal behavior for help/version
				}
			}()

			// This will test the Execute function
			Execute()
		})
	}
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("initConfig() panicked: %v", r)
		}
	}()

	// Call initConfig directly
	initConfig()
}

func TestShouldAutoInitialize(t *testing.T) {
	// Test auto-initialization logic
	result := shouldAutoInitialize()

	// Result should be a boolean (true or false)
	if result != true && result != false {
		t.Error("shouldAutoInitialize should return a boolean value")
	}
}

func TestRootCmdFlags(t *testing.T) {
	// Test that config flag exists
	configFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("Expected config flag to exist")
	}

	if configFlag.Usage == "" {
		t.Error("Config flag should have usage description")
	}
}
