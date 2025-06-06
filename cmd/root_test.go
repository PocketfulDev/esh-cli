package cmd

import (
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
