package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLastTagCmdCreation(t *testing.T) {
	if lastTagCmd == nil {
		t.Error("lastTagCmd should not be nil")
	}

	if lastTagCmd.Use != "last-tag [environment]" {
		t.Errorf("Expected Use to be 'last-tag [environment]', got %q", lastTagCmd.Use)
	}

	if lastTagCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func TestLastTagFlags(t *testing.T) {
	// Test that service flag exists
	flag := lastTagCmd.Flags().Lookup("service")
	if flag == nil {
		t.Error("Expected service flag to exist")
	}

	if flag.Shorthand != "s" {
		t.Errorf("Expected service flag shorthand to be 's', got %q", flag.Shorthand)
	}
}

func TestLastTagValidation(t *testing.T) {
	// Test argument validation
	if lastTagCmd.Args == nil {
		t.Error("Args validator should be set")
	}

	// Test with correct number of args
	err := lastTagCmd.Args(lastTagCmd, []string{"dev"})
	if err != nil {
		t.Errorf("Expected no error with 1 argument, got %v", err)
	}

	// Test with incorrect number of args
	err = lastTagCmd.Args(lastTagCmd, []string{})
	if err == nil {
		t.Error("Expected error with 0 arguments")
	}

	err = lastTagCmd.Args(lastTagCmd, []string{"dev", "extra"})
	if err == nil {
		t.Error("Expected error with 2 arguments")
	}
}

func TestLastTagCmdInRootCmd(t *testing.T) {
	// Check that last-tag command is added to root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "last-tag" {
			found = true
			break
		}
	}

	if !found {
		t.Error("last-tag command should be added to root command")
	}
}

func TestLastTagCmdHelp(t *testing.T) {
	// Test that help can be generated without error
	cmd := &cobra.Command{}
	cmd.AddCommand(lastTagCmd)

	// This should not panic
	help := lastTagCmd.Long
	if help == "" {
		t.Error("Long description should not be empty")
	}

	// Check examples exist
	if lastTagCmd.Example == "" {
		t.Error("Examples should not be empty")
	}
}
