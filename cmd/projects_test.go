package cmd

import (
	"bytes"
	"testing"
)

func TestProjectsCmdCreation(t *testing.T) {
	if projectsCmd == nil {
		t.Error("projectsCmd should not be nil")
	}

	if projectsCmd.Use != "projects" {
		t.Errorf("Expected Use to be 'projects', got %q", projectsCmd.Use)
	}

	if projectsCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func TestProjectsCmdInRootCmd(t *testing.T) {
	// Check that projects command is added to root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "projects" {
			found = true
			break
		}
	}

	if !found {
		t.Error("projects command should be added to root command")
	}
}

func TestRunProjectsExecution(t *testing.T) {
	// Test that the command can execute without panicking
	cmd := projectsCmd

	// Capture output for testing
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// This should not panic even if no config exists
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runProjects panicked: %v", r)
		}
	}()

	// Execute the command - this may produce output or error but shouldn't panic
	runProjects(cmd, []string{})
}

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		key      string
		expected string
	}{
		{"string value", map[string]interface{}{"test": "value"}, "test", "value"},
		{"missing key", map[string]interface{}{"other": "value"}, "test", "unknown"},
		{"non-string value", map[string]interface{}{"test": 123}, "test", "unknown"},
		{"nil map", nil, "test", "unknown"},
		{"empty map", map[string]interface{}{}, "test", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringValue(tt.m, tt.key)
			if result != tt.expected {
				t.Errorf("getStringValue(%v, %q) = %q, want %q", tt.m, tt.key, result, tt.expected)
			}
		})
	}
}

func TestProjectsCmdHelp(t *testing.T) {
	// Test that help can be generated without error
	if projectsCmd.Long == "" {
		t.Error("Long description should not be empty")
	}

	// Check that command has examples
	if projectsCmd.Example == "" {
		t.Error("Examples should not be empty")
	}
}
