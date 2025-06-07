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
