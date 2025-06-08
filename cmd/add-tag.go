package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	promoteFrom string
	hotFix      bool
	service     string
)

// addTagCmd represents the add-tag command
var addTagCmd = &cobra.Command{
	Use:   "add-tag [environment] [version]",
	Short: "Adds and pushes new hot fix tag",
	Long: `Adds and pushes new hot fix tag.
Tag format is env_major.minor.patch-release
In some projects this triggers deployment in CircleCI.

Use 'esh-cli last-tag [environment]' to see the current last tag.`,
	Example: `  esh-cli add-tag stg6 1.2.1 - adds tag for staging on latest commit in current directory
  esh-cli add-tag production2 1.2.1 --from stg6_1.2.1-0 - promotes from staging
  esh-cli add-tag stg6 1.2.1 --service myservice - adds tag with service prefix`,
	Args: cobra.ExactArgs(2),
	Run:  runAddTag,
}

func init() {
	rootCmd.AddCommand(addTagCmd)

	addTagCmd.Flags().StringVarP(&promoteFrom, "from", "f", "", "tag to promote from")
	addTagCmd.Flags().BoolVar(&hotFix, "hot-fix", false, "tag hot fix")
	addTagCmd.Flags().StringVarP(&service, "service", "s", "", "service name to tag")
}

func runAddTag(cmd *cobra.Command, args []string) {
	environment := args[0]
	version := args[1]

	// Validate environment
	if !utils.ContainsString(utils.ENVS, environment) {
		fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
			environment, utils.ENVS)
		os.Exit(1)
	}

	// Validate version
	if !utils.IsVersionValid(version, false) {
		fmt.Fprintf(os.Stderr, "Error: version '%s' is not valid\n", version)
		os.Exit(1)
	}

	// Validate promote_from tag if provided
	if promoteFrom != "" && !hotFix && !utils.IsTagValid(promoteFrom) {
		fmt.Fprintf(os.Stderr, "Error: tag '%s' is not valid\n", promoteFrom)
		os.Exit(1)
	}

	// Check current branch
	branch, err := utils.Cmd("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
		os.Exit(1)
	}

	// Validate branch and hot fix rules
	if version != "last" && utils.IsReleaseBranch(branch) && !hotFix {
		fmt.Fprintf(os.Stderr, "Error: you can tag only hot fix (use --hot-fix flag) from release branch\n")
		os.Exit(1)
	}

	// Check if not on master/main and not hot fix
	if branch != "master" && branch != "main" && !hotFix {
		if utils.Ask(fmt.Sprintf("Current branch is %s. Continue? (y/n)", branch)) != "y" {
			os.Exit(0)
		}
	}

	// Hot fix must be from release branch
	if hotFix && !utils.IsReleaseBranch(branch) {
		fmt.Fprintf(os.Stderr, "Error: hot fix must be tagged from release branch\n")
		os.Exit(1)
	}

	// Check if local and remote are synced
	sha, err := utils.Cmd("git rev-list HEAD | head -1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting local SHA: %v\n", err)
		os.Exit(1)
	}

	shaRemote, err := utils.Cmd(fmt.Sprintf("git rev-list origin/%s | head -1", branch))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting remote SHA: %v\n", err)
		os.Exit(1)
	}

	if sha != shaRemote {
		fmt.Fprintf(os.Stderr, "Error: remote is not synced\n")
		os.Exit(1)
	}

	// Get last tag for version
	var lastTag string
	var projectPath string

	if service == "" {
		// Use current working directory when no service specified
		projectPath = "."
		lastTag, _, err = utils.FindLastTagAndCommentInDir(environment, version, "", projectPath)
	} else {
		// Make sure config is loaded when service is specified
		initConfig()

		// Find the project path for the specified service
		projectPath = findProjectPath(service)
		if projectPath == "" {
			fmt.Fprintf(os.Stderr, "Error: service '%s' not found in configuration.\n", service)
			os.Exit(1)
		}

		// Use the project directory when service is specified
		lastTag, _, err = utils.FindLastTagAndCommentInDir(environment, version, "", projectPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding last tag in %s: %v\n", projectPath, err)
		os.Exit(1)
	}

	var newTag string
	var newTagCommit string

	if promoteFrom != "" {
		// Promote from another tag
		promoteFromEnv, err := utils.GetEnvFromTag(promoteFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing promote-from tag: %v\n", err)
			os.Exit(1)
		}

		if !utils.ContainsString(utils.ENVS, promoteFromEnv) {
			fmt.Fprintf(os.Stderr, "Error: environment '%s' does not exist\n", promoteFromEnv)
			os.Exit(1)
		}

		newTagCommit, err = utils.Cmd(fmt.Sprintf("git rev-list -n 1 %s", promoteFrom))
		if err != nil || newTagCommit == "" {
			fmt.Fprintf(os.Stderr, "Error: tag '%s' not found\n", promoteFrom)
			os.Exit(1)
		}

		newTag = strings.Replace(promoteFrom, promoteFromEnv, environment, 1)
		if utils.Ask(fmt.Sprintf("promote %s to %s? (y/n)", promoteFrom, newTag)) != "y" {
			os.Exit(0)
		}
	} else {
		// Create new tag
		newTagCommit = sha
		if lastTag != "" {
			newTag = utils.IncrementTag(lastTag, hotFix)
			if newTag == "" {
				fmt.Fprintf(os.Stderr, "Error: failed to increment tag '%s'\n", lastTag)
				os.Exit(1)
			}
		} else {
			newTag = fmt.Sprintf("%s-0", utils.TagPrefix(environment, version, service))
		}

		if utils.Ask(fmt.Sprintf("add %s? (y/n)", newTag)) != "y" {
			os.Exit(0)
		}
	}

	// Get comment for the tag
	newTagComment := utils.Ask("comment")
	if strings.TrimSpace(newTagComment) == "" {
		newTagComment = newTag
	}

	// Tag and push
	_, err = utils.Cmd(fmt.Sprintf("git tag -a %s -m \"%s\" %s", newTag, newTagComment, newTagCommit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tag: %v\n", err)
		os.Exit(1)
	}

	_, err = utils.Cmd(fmt.Sprintf("git push origin %s", newTag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pushing tag: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created and pushed tag: %s\n", newTag)
}
