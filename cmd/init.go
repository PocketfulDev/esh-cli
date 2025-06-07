package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initSearchDepth int
	initForce       bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ESH CLI with AI project discovery for pocketful and bank_1",
	Long: `Initialize ESH CLI configuration by automatically discovering projects containing "pocketful" or "bank_1".
This command will search for projects with these specific patterns and save their paths to your configuration file.`,
	Example: `  esh-cli init - discover projects containing "pocketful" or "bank_1"
  esh-cli init --depth 3 - search up to 3 directories deep
  esh-cli init --force - overwrite existing configuration`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().IntVarP(&initSearchDepth, "depth", "d", 2, "maximum search depth for project discovery")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "force overwrite existing configuration")
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🤖 ESH CLI AI Initialization Starting...")

	// Search for specific patterns: "pocketful" and "bank_1"
	targetPatterns := []string{"pocketful", "bank_1"}
	fmt.Printf("Searching for projects containing patterns: %v\n", targetPatterns)

	// Check if config already exists
	if !initForce && configExists() {
		fmt.Println("⚠️  Configuration already exists. Use --force to overwrite.")
		fmt.Printf("Current config file: %s\n", viper.ConfigFileUsed())
		return
	}

	// Discover projects with specific patterns
	projects, err := discoverSpecificProjects(targetPatterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering projects: %v\n", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Println("❌ No projects found matching the specified patterns.")
		return
	}

	// Display discovered projects
	fmt.Printf("\n🎯 Discovered %d projects:\n", len(projects))
	for i, project := range projects {
		fmt.Printf("  %d. %s (%s)\n", i+1, project.Name, project.Path)
	}

	// Save to config
	err = saveProjectsToConfig(projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Configuration saved successfully!\n")
	fmt.Printf("Config file: %s\n", getConfigFilePath())
	fmt.Printf("Found %d projects ready for ESH CLI management.\n", len(projects))
}

// Project represents a discovered project
type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// discoverProjects searches for projects on the system
func discoverProjects(projectPattern string) ([]Project, error) {
	var projects []Project

	// Get common search paths
	searchPaths := getSearchPaths()

	// Project indicators - files/folders that indicate a project
	projectIndicators := []string{
		".git",
		"package.json",
		"go.mod",
		"requirements.txt",
		"Dockerfile",
		".env",
		"docker-compose.yml",
		"pom.xml",
		"build.gradle",
		"Cargo.toml",
	}

	// Common project name patterns that might be relevant for ESH CLI
	projectNamePatterns := []string{
		"backoffice",
		"backend",
		"frontend",
		"api",
		"service",
		"app",
		"web",
		"server",
		"client",
		"admin",
		"dashboard",
		"portal",
		"platform",
		"core",
		"main",
		"primary",
	}

	for _, searchPath := range searchPaths {
		fmt.Printf("🔍 Searching in: %s\n", searchPath)

		err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // Continue on errors
			}

			// Skip if we've exceeded search depth
			if getDirectoryDepth(searchPath, path) > initSearchDepth {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip hidden directories and common ignore patterns
			if d.IsDir() && shouldSkipDirectory(d.Name()) {
				return filepath.SkipDir
			}

			// Check if this directory contains project indicators
			if d.IsDir() {
				if isProjectDirectory(path, projectIndicators) {
					projectName := filepath.Base(path)

					// Check if project name matches our patterns or contains relevant keywords
					if isRelevantProject(projectName, projectNamePatterns, projectPattern) {
						projectType := determineProjectType(path)
						projects = append(projects, Project{
							Name: projectName,
							Path: path,
							Type: projectType,
						})
					}
				}
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Warning: Error searching %s: %v\n", searchPath, err)
		}
	}

	return removeDuplicateProjects(projects), nil
}

// getSearchPaths returns common paths to search for projects
func getSearchPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return []string{
		filepath.Join(home, "WorkSpace"),
		filepath.Join(home, "workspace"),
		filepath.Join(home, "Projects"),
		filepath.Join(home, "projects"),
		filepath.Join(home, "Development"),
		filepath.Join(home, "dev"),
		filepath.Join(home, "Code"),
		filepath.Join(home, "code"),
		filepath.Join(home, "git"),
		filepath.Join(home, "src"),
		filepath.Join(home, "Documents"),
		home, // Home directory itself
		".",  // Current directory
	}
}

// shouldSkipDirectory checks if a directory should be skipped during search
func shouldSkipDirectory(name string) bool {
	skipPatterns := []string{
		".", "node_modules", "vendor", "target", "build", "dist",
		".git", ".svn", ".hg", "__pycache__", ".pytest_cache",
		"venv", "env", ".venv", ".env", "virtualenv",
		".idea", ".vscode", ".vs", ".gradle", ".m2",
		"tmp", "temp", "cache", "logs", "log",
	}

	for _, pattern := range skipPatterns {
		if strings.HasPrefix(name, pattern) {
			return true
		}
	}
	return false
}

// isProjectDirectory checks if a directory contains project indicators
func isProjectDirectory(dirPath string, indicators []string) bool {
	for _, indicator := range indicators {
		indicatorPath := filepath.Join(dirPath, indicator)
		if _, err := os.Stat(indicatorPath); err == nil {
			return true
		}
	}
	return false
}

// isRelevantProject checks if project name is relevant for ESH CLI
func isRelevantProject(projectName string, patterns []string, projectPattern string) bool {
	lowerName := strings.ToLower(projectName)

	// First, check for the specific project pattern (highest priority)
	if projectPattern != "" && strings.Contains(lowerName, strings.ToLower(projectPattern)) {
		return true
	}

	// Check for exact matches or contains patterns
	for _, pattern := range patterns {
		if strings.Contains(lowerName, pattern) {
			return true
		}
	}

	// Additional heuristics - check for common naming conventions
	if strings.Contains(lowerName, "service") ||
		strings.Contains(lowerName, "microservice") {
		return true
	}

	return false
}

// determineProjectType analyzes project directory to determine type
func determineProjectType(projectPath string) string {
	// Check for specific files to determine project type
	checks := map[string]string{
		"package.json":     "nodejs",
		"go.mod":           "golang",
		"requirements.txt": "python",
		"Cargo.toml":       "rust",
		"pom.xml":          "java",
		"build.gradle":     "gradle",
		"Dockerfile":       "docker",
	}

	for file, projectType := range checks {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return projectType
		}
	}

	return "unknown"
}

// getDirectoryDepth calculates directory depth relative to base path
func getDirectoryDepth(basePath, currentPath string) int {
	relPath, err := filepath.Rel(basePath, currentPath)
	if err != nil {
		return 0
	}

	if relPath == "." {
		return 0
	}

	return len(strings.Split(relPath, string(filepath.Separator)))
}

// removeDuplicateProjects removes duplicate projects based on path
func removeDuplicateProjects(projects []Project) []Project {
	seen := make(map[string]bool)
	var unique []Project

	for _, project := range projects {
		if !seen[project.Path] {
			seen[project.Path] = true
			unique = append(unique, project)
		}
	}

	return unique
}

// configExists checks if configuration file already exists
func configExists() bool {
	configPath := getConfigFilePath()
	_, err := os.Stat(configPath)
	return err == nil
}

// getConfigFilePath returns the path to the config file
func getConfigFilePath() string {
	if cfgFile != "" {
		return cfgFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ".esh-cli.yaml"
	}

	return filepath.Join(home, ".esh-cli.yaml")
}

// saveProjectsToConfig saves discovered projects to viper configuration
func saveProjectsToConfig(projects []Project) error {
	// Extract project paths for easy access
	var projectPaths []string
	for _, project := range projects {
		projectPaths = append(projectPaths, project.Path)
	}

	// Set up configuration structure
	config := map[string]interface{}{
		"projects":        projects,
		"project_paths":   projectPaths,
		"initialized_at":  fmt.Sprintf("%v", utils.GetCurrentTime()),
		"version":         version,
		"auto_discovered": true,
		"search_patterns": []string{"pocketful", "bank_1"},
	}

	// Set all values in viper
	for key, value := range config {
		viper.Set(key, value)
	}

	// Write config file
	configPath := getConfigFilePath()

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write configuration
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// discoverSpecificProjects searches for projects containing specific patterns
func discoverSpecificProjects(targetPatterns []string) ([]Project, error) {
	var projects []Project

	// Get common search paths
	searchPaths := getSearchPaths()

	// Project indicators - files/folders that indicate a project
	projectIndicators := []string{
		".git",
		"package.json",
		"go.mod",
		"requirements.txt",
		"Dockerfile",
		".env",
		"docker-compose.yml",
		"pom.xml",
		"build.gradle",
		"Cargo.toml",
	}

	for _, searchPath := range searchPaths {
		fmt.Printf("🔍 Searching in: %s\n", searchPath)

		err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // Continue on errors
			}

			// Skip if we've exceeded search depth
			if getDirectoryDepth(searchPath, path) > initSearchDepth {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip hidden directories and common ignore patterns
			if d.IsDir() && shouldSkipDirectory(d.Name()) {
				return filepath.SkipDir
			}

			// Check if this directory contains project indicators
			if d.IsDir() {
				if isProjectDirectory(path, projectIndicators) {
					projectName := filepath.Base(path)

					// Check if project name contains any of our target patterns
					if containsAnyPattern(projectName, targetPatterns) {
						projectType := determineProjectType(path)
						projects = append(projects, Project{
							Name: projectName,
							Path: path,
							Type: projectType,
						})
					}
				}
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Warning: Error searching %s: %v\n", searchPath, err)
		}
	}

	return removeDuplicateProjects(projects), nil
}

// containsAnyPattern checks if project name contains any of the target patterns
func containsAnyPattern(projectName string, patterns []string) bool {
	lowerName := strings.ToLower(projectName)

	for _, pattern := range patterns {
		if strings.Contains(lowerName, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}
