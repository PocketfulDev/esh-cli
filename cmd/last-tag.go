package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	Example: `  esh-cli last-tag stg6 - shows last tag for staging
  esh-cli last-tag production2 - shows last tag for production
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

	// Get last tag for environment
	lastTag, lastComment, err := utils.FindLastTagAndComment(environment, "?", lastTagService)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding last tag: %v\n", err)
		os.Exit(1)
	}

	if lastTag != "" {
		fmt.Printf("%s %s\n", lastTag, lastComment)
	} else {
		fmt.Println("No tags found")
	}
}
