package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	bumpMajor   bool
	bumpMinor   bool
	bumpPatch   bool
	bumpAuto    bool
	bumpPreview bool
	bumpService string
	fromCommit  string
)

// bumpVersionCmd represents the bump-version command
var bumpVersionCmd = &cobra.Command{
	Use:   "bump-version [environment]",
	Short: "Bump semantic version of tags",
	Long: `Bump semantic version (MAJOR.MINOR.PATCH) of tags for the specified environment.

This command creates a new tag with an incremented semantic version based on the bump type:
- --major: Increment major version (breaking changes)
- --minor: Increment minor version (new features)  
- --patch: Increment patch version (bug fixes)
- --auto: Auto-detect bump type from commit messages (conventional commits)

The new tag will have the format: env_major.minor.patch-1`,
	Example: `  esh-cli bump-version stg6 --major     # 1.2.3 → 2.0.0-1
  esh-cli bump-version stg6 --minor     # 1.2.3 → 1.3.0-1
  esh-cli bump-version stg6 --patch     # 1.2.3 → 1.2.4-1
  esh-cli bump-version stg6 --auto      # Auto-detect from commits
  esh-cli bump-version stg6 --major --preview  # Show what would be created
  esh-cli bump-version stg6 --patch --service myservice  # Service-specific tag`,
	Args: cobra.ExactArgs(1),
	Run:  runBumpVersion,
}

func init() {
	rootCmd.AddCommand(bumpVersionCmd)

	bumpVersionCmd.Flags().BoolVar(&bumpMajor, "major", false, "bump major version (breaking changes)")
	bumpVersionCmd.Flags().BoolVar(&bumpMinor, "minor", false, "bump minor version (new features)")
	bumpVersionCmd.Flags().BoolVar(&bumpPatch, "patch", false, "bump patch version (bug fixes)")
	bumpVersionCmd.Flags().BoolVar(&bumpAuto, "auto", false, "auto-detect bump type from commit messages")
	bumpVersionCmd.Flags().BoolVar(&bumpPreview, "preview", false, "preview the change without creating tag")
	bumpVersionCmd.Flags().StringVarP(&bumpService, "service", "s", "", "service name to tag")
	bumpVersionCmd.Flags().StringVar(&fromCommit, "from-commit", "HEAD", "commit to tag (default: HEAD)")

	// Mark flags as mutually exclusive
	bumpVersionCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "auto")
}

func runBumpVersion(cmd *cobra.Command, args []string) {
	environment := args[0]

	// Validate environment
	if !utils.ContainsString(utils.ENVS, environment) {
		fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
			environment, utils.ENVS)
		os.Exit(1)
	}

	// Ensure exactly one bump type is specified
	bumpCount := 0
	if bumpMajor {
		bumpCount++
	}
	if bumpMinor {
		bumpCount++
	}
	if bumpPatch {
		bumpCount++
	}
	if bumpAuto {
		bumpCount++
	}

	if bumpCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: must specify one of --major, --minor, --patch, or --auto\n")
		os.Exit(1)
	}

	// Get current working directory for tag operations
	var projectPath string
	if bumpService == "" {
		projectPath = "."
	} else {
		initConfig()
		projectPath = findProjectPath(bumpService)
		if projectPath == "" {
			fmt.Fprintf(os.Stderr, "Error: service '%s' not found in configuration.\n", bumpService)
			os.Exit(1)
		}
	}

	// Find the latest tag for the environment
	latestTag, latestVersion, err := utils.GetLatestSemanticVersion(environment, bumpService)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding latest version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current latest tag: %s (version: %s)\n", latestTag, latestVersion)

	// Determine bump type
	var bumpType utils.BumpType
	if bumpMajor {
		bumpType = utils.BumpMajor
	} else if bumpMinor {
		bumpType = utils.BumpMinor
	} else if bumpPatch {
		bumpType = utils.BumpPatch
	} else if bumpAuto {
		// Auto-detect from commits since last tag
		commits, err := utils.GetCommitsBetweenTags(latestTag, fromCommit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commits since last tag: %v\n", err)
			os.Exit(1)
		}

		if len(commits) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no commits found since last tag %s\n", latestTag)
			os.Exit(1)
		}

		bumpType = utils.DetectBumpType(commits)
		fmt.Printf("Auto-detected bump type: %s (analyzed %d commits)\n", bumpType, len(commits))
	}

	// Create new tag with bumped version
	newTag, err := utils.BumpTagVersion(latestTag, bumpType, environment, bumpService)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new tag: %v\n", err)
		os.Exit(1)
	}

	// Preview mode - show what would be created
	if bumpPreview {
		fmt.Printf("\n🔍 Preview Mode:\n")
		fmt.Printf("Current tag: %s\n", latestTag)
		fmt.Printf("Bump type:   %s\n", bumpType)
		fmt.Printf("New tag:     %s\n", newTag)
		fmt.Printf("Target commit: %s\n", fromCommit)
		fmt.Printf("\nTo create this tag, run the same command without --preview\n")
		return
	}

	// Confirm with user
	if utils.Ask(fmt.Sprintf("Create new tag %s? (y/n)", newTag)) != "y" {
		fmt.Println("Operation cancelled")
		os.Exit(0)
	}

	// Resolve target commit
	targetCommit, err := utils.Cmd(fmt.Sprintf("git rev-parse %s", fromCommit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving commit %s: %v\n", fromCommit, err)
		os.Exit(1)
	}

	// Get comment for the tag
	defaultComment := fmt.Sprintf("Bump %s version: %s", bumpType, newTag)
	comment := utils.Ask(fmt.Sprintf("Tag comment (default: %s)", defaultComment))
	if strings.TrimSpace(comment) == "" {
		comment = defaultComment
	}

	// Create and push the tag
	fmt.Printf("Creating tag %s on commit %s...\n", newTag, targetCommit[:8])

	_, err = utils.Cmd(fmt.Sprintf("git tag -a %s -m \"%s\" %s", newTag, comment, targetCommit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tag: %v\n", err)
		os.Exit(1)
	}

	_, err = utils.Cmd(fmt.Sprintf("git push origin %s", newTag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pushing tag: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully created and pushed tag: %s\n", newTag)

	// Show summary
	fmt.Printf("\n📋 Summary:\n")
	fmt.Printf("Previous: %s (%s)\n", latestTag, latestVersion)
	newVersion, _ := utils.GetVersionFromTag(newTag)
	fmt.Printf("New:      %s (%s)\n", newTag, newVersion)
	fmt.Printf("Bump:     %s\n", bumpType)
	fmt.Printf("Commit:   %s\n", targetCommit[:8])
}
