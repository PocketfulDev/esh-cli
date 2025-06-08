package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	lastTagService string
)

// lastTagCmd represents the last-tag command
var lastTagCmd = &cobra.Command{
	Use:   "last-tag [environment]",
	Short: "Shows the last tag for a given environment",
	Long: `Shows the last tag and its comment for a given environment.
This is useful for checking the current state before creating new tags.`,
	Example: `  esh-cli last-tag stg6 - shows last tag for staging in current directory
  esh-cli last-tag production2 - shows last tag for production in current directory
  esh-cli last-tag stg6 --service myservice - shows last tag for specific service`,
	Args: cobra.ExactArgs(1),
	Run:  runLastTag,
}

func init() {
	rootCmd.AddCommand(lastTagCmd)
	lastTagCmd.Flags().StringVarP(&lastTagService, "service", "s", "", "service name to check")
}

func runLastTag(cmd *cobra.Command, args []string) {
	environment := args[0]

	// Validate environment
	if !utils.ContainsString(utils.ENVS, environment) {
		fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
			environment, utils.ENVS)
		os.Exit(1)
	}

	var projectPath string

	// If no service specified, use current working directory
	if lastTagService == "" {
		projectPath = "." // Current working directory
	} else {
		// Make sure config is loaded when service is specified
		initConfig()

		// Find the project path for the specified service
		projectPath = findProjectPath(lastTagService)
		if projectPath == "" {
			fmt.Fprintf(os.Stderr, "Error: service '%s' not found in configuration.\n", lastTagService)
			fmt.Fprintf(os.Stderr, "Available services:\n")
			suggestProjects()
			os.Exit(1)
		}
	}

	// Get last tag for environment from the specific project directory (or current directory)
	// Note: We don't include the service name in the tag pattern since tags are in format: env_version-release
	lastTag, lastComment, err := utils.FindLastTagAndCommentInDir(environment, "?", "", projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding last tag in %s: %v\n", projectPath, err)
		os.Exit(1)
	}

	if lastTag != "" {
		fmt.Printf("%s %s\n", lastTag, lastComment)
	} else {
		if lastTagService == "" {
			fmt.Printf("No tags found in current directory for environment '%s'\n", environment)
		} else {
			fmt.Printf("No tags found for service '%s' in environment '%s'\n", lastTagService, environment)
		}
	}
}

// findProjectPath finds the path for a given service name from the config
func findProjectPath(serviceName string) string {
	projects := viper.Get("projects")
	if projects == nil {
		return ""
	}

	projectsList, ok := projects.([]interface{})
	if !ok {
		return ""
	}

	for _, proj := range projectsList {
		projMap, ok := proj.(map[string]interface{})
		if !ok {
			continue
		}

		name := getProjectStringValue(projMap, "name")
		path := getProjectStringValue(projMap, "path")

		// Match service name with project name
		if strings.EqualFold(name, serviceName) {
			return path
		}
	}

	return ""
}

// suggestProjects shows available projects to the user
func suggestProjects() {
	projects := viper.Get("projects")
	if projects == nil {
		fmt.Println("❌ No projects found in configuration.")
		fmt.Println("Run 'esh-cli init' to discover projects automatically.")
		return
	}

	projectsList, ok := projects.([]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Invalid projects data in configuration\n")
		return
	}

	if len(projectsList) == 0 {
		fmt.Println("❌ No projects found in configuration.")
		fmt.Println("Run 'esh-cli init' to discover projects automatically.")
		return
	}

	fmt.Println("📁 Available services/projects:")
	fmt.Println("Please specify a service using --service or -s flag:")
	fmt.Println()

	for _, proj := range projectsList {
		projMap, ok := proj.(map[string]interface{})
		if !ok {
			continue
		}

		name := getProjectStringValue(projMap, "name")
		projectType := getProjectStringValue(projMap, "type")

		fmt.Printf("  • %s (%s)\n", name, projectType)
	}

	fmt.Printf("\nExample: esh-cli last-tag %s --service %s\n",
		"stg6", getProjectStringValue(projectsList[0].(map[string]interface{}), "name"))
}

// getProjectStringValue safely extracts string value from map
func getProjectStringValue(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "unknown"
}
