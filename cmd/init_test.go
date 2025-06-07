package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmdCreation(t *testing.T) {
	if initCmd == nil {
		t.Error("initCmd should not be nil")
	}

	if initCmd.Use != "init" {
		t.Errorf("Expected Use to be 'init', got %q", initCmd.Use)
	}

	if initCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func TestInitFlags(t *testing.T) {
	// Test that depth flag exists
	depthFlag := initCmd.Flags().Lookup("depth")
	if depthFlag == nil {
		t.Error("Expected depth flag to exist")
	}

	if depthFlag.Shorthand != "d" {
		t.Errorf("Expected depth flag shorthand to be 'd', got %q", depthFlag.Shorthand)
	}

	// Test that force flag exists
	forceFlag := initCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("Expected force flag to exist")
	}

	if forceFlag.Shorthand != "f" {
		t.Errorf("Expected force flag shorthand to be 'f', got %q", forceFlag.Shorthand)
	}
}

func TestInitCmdInRootCmd(t *testing.T) {
	// Check that init command is added to root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "init" {
			found = true
			break
		}
	}

	if !found {
		t.Error("init command should be added to root command")
	}
}

func TestGetSearchPaths(t *testing.T) {
	paths := getSearchPaths()

	if len(paths) == 0 {
		t.Error("getSearchPaths() should return at least one path")
	}

	// Should contain current directory
	found := false
	for _, path := range paths {
		if path == "." {
			found = true
			break
		}
	}

	if !found {
		t.Error("getSearchPaths() should include current directory '.'")
	}
}

func TestShouldSkipDirectory(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"node_modules", true},
		{".git", true},
		{"vendor", true},
		{"normal_dir", false},
		{"my_project", false},
		{".hidden", true},
	}

	for _, tt := range tests {
		got := shouldSkipDirectory(tt.name)
		if got != tt.want {
			t.Errorf("shouldSkipDirectory(%q) = %t, want %t", tt.name, got, tt.want)
		}
	}
}

func TestIsProjectDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_project")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	indicators := []string{".git", "package.json", "go.mod"}

	// Test directory without indicators
	if isProjectDirectory(tempDir, indicators) {
		t.Error("Empty directory should not be considered a project")
	}

	// Create a .git file
	gitFile := filepath.Join(tempDir, ".git")
	if err := os.WriteFile(gitFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create .git file: %v", err)
	}

	// Test directory with indicators
	if !isProjectDirectory(tempDir, indicators) {
		t.Error("Directory with .git should be considered a project")
	}
}

func TestContainsAnyPattern(t *testing.T) {
	tests := []struct {
		projectName string
		patterns    []string
		want        bool
	}{
		{"pocketfulapp", []string{"pocketful", "bank_1"}, true},
		{"bank_1_service", []string{"pocketful", "bank_1"}, true},
		{"myapp", []string{"pocketful", "bank_1"}, false},
		{"POCKETFUL_APP", []string{"pocketful", "bank_1"}, true}, // case insensitive
		{"Bank_1_Test", []string{"pocketful", "bank_1"}, true},   // case insensitive
	}

	for _, tt := range tests {
		got := containsAnyPattern(tt.projectName, tt.patterns)
		if got != tt.want {
			t.Errorf("containsAnyPattern(%q, %v) = %t, want %t", tt.projectName, tt.patterns, got, tt.want)
		}
	}
}

func TestDetermineProjectType(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_project_type")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test unknown type (no specific files)
	projectType := determineProjectType(tempDir)
	if projectType != "unknown" {
		t.Errorf("Expected 'unknown' for empty directory, got %q", projectType)
	}

	// Test nodejs type
	packageJSON := filepath.Join(tempDir, "package.json")
	if err := os.WriteFile(packageJSON, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	projectType = determineProjectType(tempDir)
	if projectType != "nodejs" {
		t.Errorf("Expected 'nodejs' for directory with package.json, got %q", projectType)
	}
}

func TestGetDirectoryDepth(t *testing.T) {
	tests := []struct {
		basePath    string
		currentPath string
		want        int
	}{
		{"/home/user", "/home/user", 0},
		{"/home/user", "/home/user/project", 1},
		{"/home/user", "/home/user/project/src", 2},
		{"/home/user", "/home/user/a/b/c/d", 4},
	}

	for _, tt := range tests {
		got := getDirectoryDepth(tt.basePath, tt.currentPath)
		if got != tt.want {
			t.Errorf("getDirectoryDepth(%q, %q) = %d, want %d", tt.basePath, tt.currentPath, got, tt.want)
		}
	}
}

func TestRemoveDuplicateProjects(t *testing.T) {
	projects := []Project{
		{Name: "project1", Path: "/path1", Type: "nodejs"},
		{Name: "project2", Path: "/path2", Type: "golang"},
		{Name: "project1", Path: "/path1", Type: "nodejs"}, // duplicate
		{Name: "project3", Path: "/path3", Type: "python"},
	}

	unique := removeDuplicateProjects(projects)

	if len(unique) != 3 {
		t.Errorf("Expected 3 unique projects, got %d", len(unique))
	}

	// Check that paths are unique
	seen := make(map[string]bool)
	for _, project := range unique {
		if seen[project.Path] {
			t.Errorf("Found duplicate path: %s", project.Path)
		}
		seen[project.Path] = true
	}
}
