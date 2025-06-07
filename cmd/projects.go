package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List discovered projects",
	Long: `List all projects that have been discovered and saved in the configuration.
This shows projects found during initialization or manually added.`,
	Example: `  esh-cli projects - list all discovered projects
  esh-cli projects --config /path/to/config.yaml - use specific config file`,
	Run: runProjects,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}

func runProjects(cmd *cobra.Command, args []string) {
	// Make sure config is loaded
	initConfig()

	// Get projects from config
	projects := viper.Get("projects")
	if projects == nil {
		fmt.Println("❌ No projects found in configuration.")
		fmt.Println("Run 'esh-cli init' to discover projects automatically.")
		return
	}

	// Type assertion to handle viper's interface{} return
	projectsList, ok := projects.([]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Invalid projects data in configuration\n")
		os.Exit(1)
	}

	if len(projectsList) == 0 {
		fmt.Println("❌ No projects found in configuration.")
		fmt.Println("Run 'esh-cli init' to discover projects automatically.")
		return
	}

	// Display header
	fmt.Printf("📁 Found %d configured projects:\n\n", len(projectsList))

	// Display each project
	for i, proj := range projectsList {
		projMap, ok := proj.(map[string]interface{})
		if !ok {
			continue
		}

		name := getStringValue(projMap, "name")
		path := getStringValue(projMap, "path")
		projectType := getStringValue(projMap, "type")

		fmt.Printf("  %d. %s\n", i+1, name)
		fmt.Printf("     Path: %s\n", path)
		fmt.Printf("     Type: %s\n", projectType)
		fmt.Println()
	}

	// Show config info
	fmt.Printf("Configuration file: %s\n", viper.ConfigFileUsed())

	// Show initialization info if available
	if initTime := viper.GetString("initialized_at"); initTime != "" {
		fmt.Printf("Initialized at: %s\n", initTime)
	}

	if autoDiscovered := viper.GetBool("auto_discovered"); autoDiscovered {
		fmt.Println("Projects were auto-discovered. Run 'esh-cli init --force' to refresh.")
	}
}

// getStringValue safely extracts string value from map
func getStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "unknown"
}
