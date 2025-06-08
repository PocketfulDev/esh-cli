package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	diffShowHistory bool
	diffShowRemote  bool
	diffShowCommits bool
	diffShowFiles   bool
	diffShowStats   bool
	diffSinceDate   string
)

// versionDiffCmd represents the version-diff command
var versionDiffCmd = &cobra.Command{
	Use:   "version-diff <tag1> [tag2] [flags]",
	Short: "Compare versions and analyze semantic differences",
	Long: `Compare two versions or tags and analyze their semantic differences.

This command helps understand the changes between versions, including:
- Semantic version differences (major, minor, patch)
- Commits between versions
- Files changed
- Release timeline analysis

If only one tag is provided, it compares with the previous version.
If no tags are provided, it analyzes the current environment.`,
	Example: `  esh-cli version-diff stg6_1.2.3-1 stg6_1.2.4-1    # Compare two specific tags
  esh-cli version-diff stg6_1.2.3-1 --commits         # Show commits since this tag
  esh-cli version-diff stg6 --history                 # Show version history for environment
  esh-cli version-diff --since 2024-01-01             # Show changes since date`,
	Args: cobra.MinimumNArgs(1),
	Run:  runVersionDiff,
}

func init() {
	rootCmd.AddCommand(versionDiffCmd)

	versionDiffCmd.Flags().BoolVar(&diffShowHistory, "history", false, "Show version history")
	versionDiffCmd.Flags().BoolVar(&diffShowRemote, "remote", false, "Compare with remote tags")
	versionDiffCmd.Flags().BoolVar(&diffShowCommits, "commits", false, "Show commits between versions")
	versionDiffCmd.Flags().BoolVar(&diffShowFiles, "files", false, "Show changed files")
	versionDiffCmd.Flags().BoolVar(&diffShowStats, "stats", false, "Show detailed statistics")
	versionDiffCmd.Flags().StringVar(&diffSinceDate, "since", "", "Show changes since date (YYYY-MM-DD)")
}

func runVersionDiff(cmd *cobra.Command, args []string) {
	if len(args) == 1 && !utils.IsTagValid(args[0]) {
		// First argument is environment, show environment history
		environment := args[0]
		if !utils.ContainsString(utils.ENVS, environment) {
			fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
				environment, utils.ENVS)
			os.Exit(1)
		}
		showEnvironmentHistory(environment)
		return
	}

	tag1 := args[0]
	var tag2 string

	if len(args) > 1 {
		tag2 = args[1]
	} else {
		// Find previous tag automatically
		var err error
		tag2, err = findPreviousTag(tag1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding previous tag: %v\n", err)
			os.Exit(1)
		}
	}

	// Validate tags
	if !utils.IsTagValid(tag1) {
		fmt.Fprintf(os.Stderr, "Error: invalid tag format '%s'\n", tag1)
		os.Exit(1)
	}
	if tag2 != "" && !utils.IsTagValid(tag2) {
		fmt.Fprintf(os.Stderr, "Error: invalid tag format '%s'\n", tag2)
		os.Exit(1)
	}

	// Compare versions
	err := compareVersions(tag1, tag2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error comparing versions: %v\n", err)
		os.Exit(1)
	}
}

func showEnvironmentHistory(environment string) {
	fmt.Printf("📊 Version History for Environment: %s\n\n", environment)

	// Get all tags for environment
	pattern := fmt.Sprintf("*_%s_*", environment)
	output, err := utils.Cmd(fmt.Sprintf("git tag -l '%s' --sort=-version:refname", pattern))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing tags: %v\n", err)
		os.Exit(1)
	}

	if output == "" {
		fmt.Printf("No tags found for environment '%s'\n", environment)
		return
	}

	tags := strings.Split(output, "\n")
	if len(tags) == 0 {
		fmt.Printf("No tags found for environment '%s'\n", environment)
		return
	}

	fmt.Printf("Found %d versions:\n\n", len(tags))

	for i, tag := range tags {
		if tag == "" {
			continue
		}

		version, err := utils.GetVersionFromTag(tag)
		if err != nil {
			fmt.Printf("⚠️  %s (invalid version format)\n", tag)
			continue
		}

		// Get tag date
		dateStr, err := utils.Cmd(fmt.Sprintf("git log -1 --format=%%ai %s", tag))
		if err != nil {
			dateStr = "unknown"
		}

		var date time.Time
		if dateStr != "unknown" {
			date, _ = time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		}

		// Show semantic version difference from previous
		diffStr := ""
		if i < len(tags)-1 && tags[i+1] != "" {
			prevVersion, err := utils.GetVersionFromTag(tags[i+1])
			if err == nil {
				bumpType := getBumpType(prevVersion, version)
				diffStr = fmt.Sprintf(" (%s)", bumpType)
			}
		}

		if date.IsZero() {
			fmt.Printf("  %s (%s)%s\n", tag, version, diffStr)
		} else {
			fmt.Printf("  %s (%s) - %s%s\n", tag, version, date.Format("2006-01-02"), diffStr)
		}
	}

	if diffShowStats {
		showEnvironmentStats(environment, tags)
	}
}

func compareVersions(tag1, tag2 string) error {
	if tag2 == "" {
		fmt.Printf("📋 Analyzing Version: %s\n\n", tag1)
	} else {
		fmt.Printf("📋 Comparing Versions: %s → %s\n\n", tag2, tag1)
	}

	// Get versions
	version1, err := utils.GetVersionFromTag(tag1)
	if err != nil {
		return fmt.Errorf("error parsing tag1 version: %v", err)
	}

	var version2 string
	if tag2 != "" {
		version2, err = utils.GetVersionFromTag(tag2)
		if err != nil {
			return fmt.Errorf("error parsing tag2 version: %v", err)
		}

		// Show semantic difference
		bumpType := getBumpType(version2, version1)
		fmt.Printf("Semantic Change: %s → %s (%s)\n", version2, version1, bumpType)
	} else {
		fmt.Printf("Version: %s\n", version1)
	}

	// Show commits if requested or if no tag2
	if diffShowCommits || tag2 == "" {
		if tag2 != "" {
			fmt.Printf("\n📝 Commits between %s and %s:\n", tag2, tag1)
			commits, err := utils.GetCommitsBetweenTags(tag2, tag1)
			if err != nil {
				return fmt.Errorf("error getting commits: %v", err)
			}
			showCommits(commits)
		} else {
			fmt.Printf("\n📝 Recent commits up to %s:\n", tag1)
			// Show last 10 commits up to tag
			output, err := utils.Cmd(fmt.Sprintf("git log --oneline -10 %s", tag1))
			if err != nil {
				return fmt.Errorf("error getting commits: %v", err)
			}
			if output != "" {
				commits := strings.Split(output, "\n")
				for _, commit := range commits {
					if commit != "" {
						fmt.Printf("  • %s\n", commit)
					}
				}
			}
		}
	}

	// Show files if requested
	if diffShowFiles && tag2 != "" {
		fmt.Printf("\n📁 Changed Files:\n")
		output, err := utils.Cmd(fmt.Sprintf("git diff --name-only %s..%s", tag2, tag1))
		if err != nil {
			return fmt.Errorf("error getting changed files: %v", err)
		}
		if output != "" {
			files := strings.Split(output, "\n")
			for _, file := range files {
				if file != "" {
					fmt.Printf("  • %s\n", file)
				}
			}
		} else {
			fmt.Println("  No files changed")
		}
	}

	// Show stats if requested
	if diffShowStats && tag2 != "" {
		fmt.Printf("\n📊 Statistics:\n")
		showDiffStats(tag2, tag1)
	}

	return nil
}

func findPreviousTag(tag string) (string, error) {
	// Extract environment from tag
	parts := strings.Split(tag, "_")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid tag format")
	}

	var env string
	if len(parts) == 3 {
		// service_env_version format
		env = parts[1]
	} else {
		// env_version format
		env = parts[0]
	}

	// Get all tags for environment, sorted by version
	pattern := fmt.Sprintf("*_%s_*", env)
	output, err := utils.Cmd(fmt.Sprintf("git tag -l '%s' --sort=-version:refname", pattern))
	if err != nil {
		return "", err
	}

	if output == "" {
		return "", fmt.Errorf("no tags found for environment")
	}

	tags := strings.Split(output, "\n")

	// Find current tag and return the next one
	for i, t := range tags {
		if t == tag && i+1 < len(tags) {
			return tags[i+1], nil
		}
	}

	return "", fmt.Errorf("previous tag not found")
}

func getBumpType(oldVersion, newVersion string) string {
	oldSV, err1 := utils.ParseSemanticVersion(oldVersion)
	newSV, err2 := utils.ParseSemanticVersion(newVersion)

	if err1 != nil || err2 != nil {
		return "unknown"
	}

	if newSV.Major > oldSV.Major {
		return "MAJOR"
	} else if newSV.Minor > oldSV.Minor {
		return "MINOR"
	} else if newSV.Patch > oldSV.Patch {
		return "PATCH"
	}

	return "none"
}

func showCommits(commits []string) {
	if len(commits) == 0 {
		fmt.Println("  No commits found")
		return
	}

	for _, commit := range commits {
		if commit != "" {
			fmt.Printf("  • %s\n", commit)
		}
	}
}

func showDiffStats(tag1, tag2 string) {
	// Get commit count
	output, err := utils.Cmd(fmt.Sprintf("git rev-list --count %s..%s", tag1, tag2))
	if err == nil && output != "" {
		fmt.Printf("  Commits: %s\n", strings.TrimSpace(output))
	}

	// Get file changes
	output, err = utils.Cmd(fmt.Sprintf("git diff --shortstat %s..%s", tag1, tag2))
	if err == nil && output != "" {
		fmt.Printf("  Changes: %s\n", strings.TrimSpace(output))
	}

	// Get contributors
	output, err = utils.Cmd(fmt.Sprintf("git shortlog -sn %s..%s", tag1, tag2))
	if err == nil && output != "" {
		lines := strings.Split(output, "\n")
		fmt.Printf("  Contributors: %d\n", len(lines))
	}
}

func showEnvironmentStats(environment string, tags []string) {
	fmt.Printf("\n📊 Environment Statistics:\n")

	majorCount := 0
	minorCount := 0
	patchCount := 0

	for i := 0; i < len(tags)-1; i++ {
		if tags[i] == "" || tags[i+1] == "" {
			continue
		}

		v1, err1 := utils.GetVersionFromTag(tags[i])
		v2, err2 := utils.GetVersionFromTag(tags[i+1])

		if err1 != nil || err2 != nil {
			continue
		}

		bumpType := getBumpType(v2, v1)
		switch bumpType {
		case "MAJOR":
			majorCount++
		case "MINOR":
			minorCount++
		case "PATCH":
			patchCount++
		}
	}

	total := majorCount + minorCount + patchCount
	fmt.Printf("  Total Releases: %d\n", len(tags))
	fmt.Printf("  Major Bumps: %d\n", majorCount)
	fmt.Printf("  Minor Bumps: %d\n", minorCount)
	fmt.Printf("  Patch Bumps: %d\n", patchCount)

	if total > 0 {
		fmt.Printf("  Release Types: %.1f%% patch, %.1f%% minor, %.1f%% major\n",
			float64(patchCount)/float64(total)*100,
			float64(minorCount)/float64(total)*100,
			float64(majorCount)/float64(total)*100)
	}
}
