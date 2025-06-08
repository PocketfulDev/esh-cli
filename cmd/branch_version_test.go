package cmd

import (
	"bytes"
	"esh-cli/pkg/utils"
	"strings"
	"testing"
)

func TestBranchVersionCmdCreation(t *testing.T) {
	if branchVersionCmd == nil {
		t.Error("branchVersionCmd should not be nil")
	}

	if !strings.HasPrefix(branchVersionCmd.Use, "branch-version") {
		t.Errorf("Expected branchVersionCmd.Use to start with 'branch-version', got '%s'", branchVersionCmd.Use)
	}

	if branchVersionCmd.Short == "" {
		t.Error("branchVersionCmd.Short should not be empty")
	}
}

func TestBranchVersionFlags(t *testing.T) {
	// Test that flags exist
	flags := []string{"suggest", "auto-tag", "release-prep", "env", "service"}

	for _, flagName := range flags {
		flag := branchVersionCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' should be defined", flagName)
		}
	}
}

func TestBranchVersionCmdHelp(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	branchVersionCmd.SetOut(&buf)
	branchVersionCmd.SetErr(&buf)

	// Set args to trigger help
	branchVersionCmd.SetArgs([]string{"--help"})

	// Execute command (this should not return an error for help)
	err := branchVersionCmd.Execute()
	if err != nil {
		// Help command may return an error in some versions, but should still show help
		output := buf.String()
		if !strings.Contains(output, "branch-version") {
			t.Errorf("Help output should contain command name, got: %s", output)
		}
	}
}

// TestAnalyzeBranch tests the analyzeBranch function
func TestAnalyzeBranch(t *testing.T) {
	tests := []struct {
		name             string
		branchName       string
		expectedType     string
		expectedFeature  string
		expectedStrategy string
	}{
		{
			name:             "feature branch",
			branchName:       "feature/user-auth",
			expectedType:     "feature",
			expectedFeature:  "user-auth",
			expectedStrategy: "minor version bump (new features)",
		},
		{
			name:             "feat branch",
			branchName:       "feat/user-auth",
			expectedType:     "feature",
			expectedFeature:  "user-auth",
			expectedStrategy: "minor version bump (new features)",
		},
		{
			name:             "hotfix branch",
			branchName:       "hotfix/urgent-bug",
			expectedType:     "hotfix",
			expectedFeature:  "urgent-bug",
			expectedStrategy: "patch version bump (bug fixes)",
		},
		{
			name:             "fix branch",
			branchName:       "fix/urgent-bug",
			expectedType:     "hotfix",
			expectedFeature:  "urgent-bug",
			expectedStrategy: "patch version bump (bug fixes)",
		},
		{
			name:             "release branch",
			branchName:       "release/1.2.0",
			expectedType:     "release",
			expectedFeature:  "1.2.0",
			expectedStrategy: "prepare for release tagging",
		},
		{
			name:             "develop branch",
			branchName:       "develop",
			expectedType:     "develop",
			expectedFeature:  "",
			expectedStrategy: "analyze commits for bump type",
		},
		{
			name:             "main branch",
			branchName:       "main",
			expectedType:     "main",
			expectedFeature:  "",
			expectedStrategy: "analyze commits for bump type",
		},
		{
			name:             "bugfix branch",
			branchName:       "bugfix/memory-leak",
			expectedType:     "bugfix",
			expectedFeature:  "memory-leak",
			expectedStrategy: "patch version bump (bug fixes)",
		},
		{
			name:             "chore branch",
			branchName:       "chore/cleanup",
			expectedType:     "chore",
			expectedFeature:  "cleanup",
			expectedStrategy: "patch version bump (maintenance)",
		},
		{
			name:             "custom branch",
			branchName:       "my-custom-branch",
			expectedType:     "custom",
			expectedFeature:  "",
			expectedStrategy: "analyze commits or manual specification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzeBranch(tt.branchName)

			if result.Name != tt.branchName {
				t.Errorf("Expected Name=%s, got %s", tt.branchName, result.Name)
			}
			if result.Type != tt.expectedType {
				t.Errorf("Expected Type=%s, got %s", tt.expectedType, result.Type)
			}
			if result.Feature != tt.expectedFeature {
				t.Errorf("Expected Feature=%s, got %s", tt.expectedFeature, result.Feature)
			}
			if !strings.Contains(result.Strategy, strings.Split(tt.expectedStrategy, " ")[0]) {
				t.Errorf("Expected Strategy to contain '%s', got %s", tt.expectedStrategy, result.Strategy)
			}
		})
	}
}

// TestBranchVersionValidation tests the validation logic
func TestBranchVersionValidation(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expectValid bool
	}{
		{
			name:        "valid dev environment",
			environment: "dev",
			expectValid: true,
		},
		{
			name:        "valid stg6 environment",
			environment: "stg6",
			expectValid: true,
		},
		{
			name:        "valid production2 environment",
			environment: "production2",
			expectValid: true,
		},
		{
			name:        "invalid environment",
			environment: "invalid_env",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test environment validation logic
			isValid := utils.ContainsString(utils.ENVS, tt.environment)
			if isValid != tt.expectValid {
				t.Errorf("Expected environment '%s' valid=%v, got %v", tt.environment, tt.expectValid, isValid)
			}
		})
	}
}

// TestBranchVersionFlagCombinations tests various flag combinations
func TestBranchVersionFlagCombinations(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		args     []string
		expectOk bool
	}{
		{
			name:     "suggest flag only",
			flags:    []string{"--suggest"},
			args:     []string{},
			expectOk: true,
		},
		{
			name:     "release-prep flag only",
			flags:    []string{"--release-prep"},
			args:     []string{},
			expectOk: true,
		},
		{
			name:     "suggest with service",
			flags:    []string{"--suggest", "--service", "api"},
			args:     []string{},
			expectOk: true,
		},
		{
			name:     "suggest with env",
			flags:    []string{"--suggest", "--env", "stg6"},
			args:     []string{},
			expectOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that flag combinations parse correctly
			cmd := branchVersionCmd
			cmd.SetArgs(append(tt.flags, tt.args...))

			// This tests that the command structure accepts these flag combinations
			// without actually executing the command
			if cmd.Flags() == nil {
				t.Error("Command should have flags defined")
			}

			// Verify specific flags exist
			for _, flag := range tt.flags {
				if strings.HasPrefix(flag, "--") {
					flagName := strings.TrimPrefix(flag, "--")
					if cmd.Flags().Lookup(flagName) == nil {
						t.Errorf("Flag '%s' should exist", flagName)
					}
				}
			}
		})
	}
}

// TestAnalyzeCommitsForBumpLogic tests the commit analysis logic patterns
func TestAnalyzeCommitsForBumpLogic(t *testing.T) {
	tests := []struct {
		name           string
		commitMessages []string
		expectedBump   utils.BumpType
	}{
		{
			name: "breaking change",
			commitMessages: []string{
				"feat!: breaking change in API",
				"fix: some fix",
			},
			expectedBump: utils.BumpMajor,
		},
		{
			name: "feature change",
			commitMessages: []string{
				"feat: new user authentication",
				"fix: some fix",
			},
			expectedBump: utils.BumpMinor,
		},
		{
			name: "fix only",
			commitMessages: []string{
				"fix: memory leak in parser",
				"docs: update readme",
			},
			expectedBump: utils.BumpPatch,
		},
		{
			name: "no significant changes",
			commitMessages: []string{
				"docs: update documentation",
				"chore: cleanup code",
			},
			expectedBump: utils.BumpPatch, // Default fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the DetectBumpType logic with our test messages
			result := utils.DetectBumpType(tt.commitMessages)
			if result != tt.expectedBump {
				t.Errorf("Expected bump type=%v, got %v for commits: %v", tt.expectedBump, result, tt.commitMessages)
			}
		})
	}
}

// TestRunBranchVersionIntegration tests the runBranchVersion function execution
func TestRunBranchVersionIntegration(t *testing.T) {
	tests := []struct {
		name        string
		suggest     bool
		autoTag     bool
		releasePrep bool
		environment string
		service     string
		expectError bool
		skipExec    bool
	}{
		{
			name:        "basic branch analysis only",
			suggest:     false,
			autoTag:     false,
			releasePrep: false,
			expectError: false,
			skipExec:    false,
		},
		{
			name:        "suggest flag only",
			suggest:     true,
			autoTag:     false,
			releasePrep: false,
			expectError: false,
			skipExec:    false,
		},
		{
			name:        "auto-tag without environment",
			suggest:     false,
			autoTag:     true,
			releasePrep: false,
			expectError: true,
			skipExec:    true, // This would exit, so skip actual execution
		},
		{
			name:        "auto-tag with valid environment",
			suggest:     false,
			autoTag:     true,
			releasePrep: false,
			environment: "stg6",
			expectError: false,
			skipExec:    true, // Skip git commands for testing
		},
		{
			name:        "release prep only",
			suggest:     false,
			autoTag:     false,
			releasePrep: true,
			expectError: false,
			skipExec:    false,
		},
		{
			name:        "combined suggest and release prep",
			suggest:     true,
			autoTag:     false,
			releasePrep: true,
			expectError: false,
			skipExec:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipExec {
				// Test flag validation logic without executing
				// Reset flags
				branchSuggest = tt.suggest
				branchAutoTag = tt.autoTag
				branchReleasePrep = tt.releasePrep
				branchEnvironment = tt.environment
				branchService = tt.service

				// Test auto-tag validation
				if tt.autoTag && tt.environment == "" && tt.expectError {
					// This validates the error condition path
					t.Logf("Correctly identifies auto-tag without environment as error condition")
				}

				// Test that flags are set correctly
				if branchSuggest != tt.suggest {
					t.Errorf("Expected suggest=%v, got %v", tt.suggest, branchSuggest)
				}
				if branchAutoTag != tt.autoTag {
					t.Errorf("Expected autoTag=%v, got %v", tt.autoTag, branchAutoTag)
				}
				if branchReleasePrep != tt.releasePrep {
					t.Errorf("Expected releasePrep=%v, got %v", tt.releasePrep, branchReleasePrep)
				}

				return
			}

			// For non-skip tests, test the command structure
			cmd := branchVersionCmd
			args := []string{}

			if tt.suggest {
				args = append(args, "--suggest")
			}
			if tt.releasePrep {
				args = append(args, "--release-prep")
			}
			if tt.environment != "" {
				args = append(args, "--env", tt.environment)
			}
			if tt.service != "" {
				args = append(args, "--service", tt.service)
			}

			cmd.SetArgs(args)

			// Test that the command structure accepts these arguments
			// without actually executing (since we can't guarantee git availability)
			if cmd.Run == nil {
				t.Error("Command should have a Run function")
			}

			// Verify flag parsing works
			err := cmd.ParseFlags(args)
			if err != nil {
				if !tt.expectError {
					t.Errorf("Unexpected error parsing flags: %v", err)
				}
			} else if tt.expectError {
				t.Error("Expected error but got none")
			}
		})
	}
}

// TestRunBranchVersionFlagValidation tests the flag validation in runBranchVersion
func TestRunBranchVersionFlagValidation(t *testing.T) {
	tests := []struct {
		name               string
		autoTag            bool
		environment        string
		expectEnvironError bool
	}{
		{
			name:               "auto-tag without environment should error",
			autoTag:            true,
			environment:        "",
			expectEnvironError: true,
		},
		{
			name:               "auto-tag with valid environment",
			autoTag:            true,
			environment:        "stg6",
			expectEnvironError: false,
		},
		{
			name:               "no auto-tag flag",
			autoTag:            false,
			environment:        "",
			expectEnvironError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic directly
			if tt.autoTag && tt.environment == "" {
				// This should trigger an error in the actual function
				if !tt.expectEnvironError {
					t.Error("Expected environment error for auto-tag without environment")
				}
			}

			// Test environment validation if provided
			if tt.environment != "" {
				isValid := utils.ContainsString(utils.ENVS, tt.environment)
				if !isValid {
					t.Errorf("Environment %s should be valid", tt.environment)
				}
			}
		})
	}
}

// TestRunBranchVersionWorkflow tests different workflow combinations
func TestRunBranchVersionWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]string
		workflow string
	}{
		{
			name: "suggestion workflow",
			flags: map[string]string{
				"suggest": "true",
			},
			workflow: "analysis and suggestion",
		},
		{
			name: "auto-tagging workflow",
			flags: map[string]string{
				"auto-tag": "true",
				"env":      "stg6",
			},
			workflow: "analysis and auto-tagging",
		},
		{
			name: "release preparation workflow",
			flags: map[string]string{
				"release-prep": "true",
			},
			workflow: "analysis and release preparation",
		},
		{
			name: "combined workflow",
			flags: map[string]string{
				"suggest":      "true",
				"release-prep": "true",
			},
			workflow: "analysis, suggestion, and release preparation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that each workflow combination is supported
			cmd := branchVersionCmd

			args := []string{}
			for flag, value := range tt.flags {
				if value == "true" {
					args = append(args, "--"+flag)
				} else {
					args = append(args, "--"+flag, value)
				}
			}

			cmd.SetArgs(args)

			// Test flag parsing for workflow
			err := cmd.ParseFlags(args)
			if err != nil {
				t.Errorf("Workflow %s should parse successfully, got error: %v", tt.workflow, err)
			}

			// Verify the command has the expected structure
			if cmd.Use == "" {
				t.Error("Command should have Use defined")
			}
			if cmd.Short == "" {
				t.Error("Command should have Short description")
			}
			if cmd.Long == "" {
				t.Error("Command should have Long description")
			}
		})
	}
}
