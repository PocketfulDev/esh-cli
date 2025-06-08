package cmd

import (
	"os"
	"os/exec"
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

func TestRunLastTag(t *testing.T) {
	// Test the runLastTag function with valid environment
	// This test is designed to achieve coverage of the runLastTag function

	tests := []struct {
		name       string
		args       []string
		service    string
		shouldExit bool
		setup      func()
		cleanup    func()
	}{
		{
			name:       "valid environment without service",
			args:       []string{"dev"},
			service:    "",
			shouldExit: false,
			setup: func() {
				lastTagService = ""
			},
			cleanup: func() {
				lastTagService = ""
			},
		},
		{
			name:       "invalid environment",
			args:       []string{"invalid"},
			service:    "",
			shouldExit: true,
			setup: func() {
				lastTagService = ""
			},
			cleanup: func() {
				lastTagService = ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()

			// Create a mock command for testing
			cmd := &cobra.Command{}

			if tt.shouldExit {
				// Test cases that should call os.Exit
				// We use a subprocess approach to handle os.Exit
				if os.Getenv("BE_CRASHER") == "1" {
					runLastTag(cmd, tt.args)
					return
				}

				// Run the test in a subprocess
				subCmd := exec.Command(os.Args[0], "-test.run=TestRunLastTag/"+tt.name)
				subCmd.Env = append(os.Environ(), "BE_CRASHER=1")
				err := subCmd.Run()
				if e, ok := err.(*exec.ExitError); ok && !e.Success() {
					// Expected os.Exit, test passed
					return
				}
				t.Errorf("Expected process to exit with error, but it didn't")
			} else {
				// Test cases that should not call os.Exit
				// We'll capture any panics or issues
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("runLastTag panicked: %v", r)
					}
				}()

				// Call the function - this should execute without calling os.Exit
				runLastTag(cmd, tt.args)
			}
		})
	}
}

func TestRunLastTagWithService(t *testing.T) {
	// Test runLastTag with service specified
	// This requires config setup to achieve proper coverage

	// Setup temporary config for testing
	originalService := lastTagService
	defer func() {
		lastTagService = originalService
	}()

	// Create a test config scenario
	viper.Set("projects", []interface{}{
		map[string]interface{}{
			"name": "test-service",
			"path": ".",
		},
	})
	defer viper.Reset()

	tests := []struct {
		name       string
		service    string
		args       []string
		shouldExit bool
	}{
		{
			name:       "valid service with valid environment",
			service:    "test-service",
			args:       []string{"dev"},
			shouldExit: false,
		},
		{
			name:       "invalid service",
			service:    "non-existent-service",
			args:       []string{"dev"},
			shouldExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastTagService = tt.service
			cmd := &cobra.Command{}

			if tt.shouldExit {
				// Test subprocess approach for os.Exit cases
				if os.Getenv("BE_CRASHER") == "1" {
					runLastTag(cmd, tt.args)
					return
				}

				subCmd := exec.Command(os.Args[0], "-test.run=TestRunLastTagWithService/"+tt.name)
				subCmd.Env = append(os.Environ(), "BE_CRASHER=1")
				err := subCmd.Run()
				if e, ok := err.(*exec.ExitError); ok && !e.Success() {
					return // Expected exit
				}
				t.Errorf("Expected process to exit with error")
			} else {
				// Non-exit case - should execute successfully
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("runLastTag panicked: %v", r)
					}
				}()
				runLastTag(cmd, tt.args)
			}
		})
	}
}

func TestFindProjectPath(t *testing.T) {
	// Test the findProjectPath helper function
	// Reset viper after test
	defer viper.Reset()

	tests := []struct {
		name        string
		serviceName string
		projects    interface{}
		expected    string
	}{
		{
			name:        "service found",
			serviceName: "test-service",
			projects: []interface{}{
				map[string]interface{}{
					"name": "test-service",
					"path": "/path/to/service",
				},
			},
			expected: "/path/to/service",
		},
		{
			name:        "service not found",
			serviceName: "missing-service",
			projects: []interface{}{
				map[string]interface{}{
					"name": "other-service",
					"path": "/path/to/other",
				},
			},
			expected: "",
		},
		{
			name:        "no projects configured",
			serviceName: "any-service",
			projects:    nil,
			expected:    "",
		},
		{
			name:        "invalid projects format",
			serviceName: "any-service",
			projects:    "invalid",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("projects", tt.projects)

			result := findProjectPath(tt.serviceName)
			if result != tt.expected {
				t.Errorf("findProjectPath(%q) = %q, expected %q", tt.serviceName, result, tt.expected)
			}
		})
	}
}

func TestSuggestProjects(t *testing.T) {
	// Test the suggestProjects function for coverage
	defer viper.Reset()

	// Test with no projects
	viper.Set("projects", nil)
	// This function prints to stderr, we just want to call it for coverage
	suggestProjects()

	// Test with projects
	viper.Set("projects", []interface{}{
		map[string]interface{}{
			"name": "service1",
			"path": "/path1",
		},
		map[string]interface{}{
			"name": "service2",
			"path": "/path2",
		},
	})
	suggestProjects()
}

func TestGetProjectStringValue(t *testing.T) {
	// Test the getProjectStringValue helper function
	tests := []struct {
		name     string
		project  map[string]interface{}
		key      string
		expected string
	}{
		{
			name: "valid map with string value",
			project: map[string]interface{}{
				"name": "test-service",
				"path": "/test/path",
			},
			key:      "name",
			expected: "test-service",
		},
		{
			name: "valid map with missing key",
			project: map[string]interface{}{
				"name": "test-service",
			},
			key:      "path",
			expected: "unknown",
		},
		{
			name: "non-string value",
			project: map[string]interface{}{
				"count": 123,
			},
			key:      "count",
			expected: "unknown",
		},
		{
			name:     "empty map",
			project:  map[string]interface{}{},
			key:      "name",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProjectStringValue(tt.project, tt.key)
			if result != tt.expected {
				t.Errorf("getProjectStringValue(%v, %q) = %q, expected %q",
					tt.project, tt.key, result, tt.expected)
			}
		})
	}
}
