package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func TestFindProjectPath(t *testing.T) {
	// Setup test configuration
	originalProjects := viper.Get("projects")
	defer func() {
		viper.Set("projects", originalProjects)
	}()

	testProjects := []interface{}{
		map[string]interface{}{
			"name": "pocketfulbackoffice",
			"path": "/Users/jonathanpick/WorkSpace/GetPocketful/pocketfulbackoffice",
			"type": "nodejs",
		},
		map[string]interface{}{
			"name": "pocketfulserver",
			"path": "/Users/jonathanpick/WorkSpace/GetPocketful/pocketfulserver",
			"type": "nodejs",
		},
	}
	viper.Set("projects", testProjects)

	tests := []struct {
		name         string
		serviceName  string
		expectedPath string
	}{
		{
			name:         "existing service",
			serviceName:  "pocketfulbackoffice",
			expectedPath: "/Users/jonathanpick/WorkSpace/GetPocketful/pocketfulbackoffice",
		},
		{
			name:         "case insensitive match",
			serviceName:  "POCKETFULSERVER",
			expectedPath: "/Users/jonathanpick/WorkSpace/GetPocketful/pocketfulserver",
		},
		{
			name:         "non-existing service",
			serviceName:  "nonexistent",
			expectedPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := findProjectPath(tt.serviceName)
			if path != tt.expectedPath {
				t.Errorf("findProjectPath(%q) = %q, want %q", tt.serviceName, path, tt.expectedPath)
			}
		})
	}
}

func TestGetProjectStringValue(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "existing string value",
			m:        map[string]interface{}{"name": "test"},
			key:      "name",
			expected: "test",
		},
		{
			name:     "non-existing key",
			m:        map[string]interface{}{"name": "test"},
			key:      "missing",
			expected: "unknown",
		},
		{
			name:     "non-string value",
			m:        map[string]interface{}{"count": 42},
			key:      "count",
			expected: "unknown",
		},
		{
			name:     "empty map",
			m:        map[string]interface{}{},
			key:      "name",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProjectStringValue(tt.m, tt.key)
			if result != tt.expected {
				t.Errorf("getProjectStringValue(%v, %q) = %q, want %q", tt.m, tt.key, result, tt.expected)
			}
		})
	}
}

func TestLastTagServiceValidation(t *testing.T) {
	// Setup test configuration
	originalProjects := viper.Get("projects")
	defer func() {
		viper.Set("projects", originalProjects)
	}()

	testProjects := []interface{}{
		map[string]interface{}{
			"name": "pocketfulbackoffice",
			"path": "/Users/jonathanpick/WorkSpace/GetPocketful/pocketfulbackoffice",
			"type": "nodejs",
		},
	}
	viper.Set("projects", testProjects)

	tests := []struct {
		name       string
		service    string
		shouldFind bool
	}{
		{
			name:       "valid service",
			service:    "pocketfulbackoffice",
			shouldFind: true,
		},
		{
			name:       "invalid service",
			service:    "invalidservice",
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := findProjectPath(tt.service)
			found := path != ""
			if found != tt.shouldFind {
				t.Errorf("findProjectPath(%q) found = %v, want %v", tt.service, found, tt.shouldFind)
			}
		})
	}
}
