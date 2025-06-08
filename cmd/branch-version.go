package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	branchSuggest     bool
	branchAutoTag     bool
	branchReleasePrep bool
	branchEnvironment string
	branchService     string
)

// branchVersionCmd represents the branch-version command
var branchVersionCmd = &cobra.Command{
	Use:   "branch-version [flags]",
	Short: "Branch-aware versioning and git flow integration",
	Long: `Branch-aware versioning and git flow integration.

This command provides intelligent version management based on the current git branch:
- Suggests appropriate version bumps based on branch naming conventions
- Supports git flow workflows (feature/, hotfix/, release/ branches)
- Can automatically create tags for release preparation
- Integrates with conventional commit analysis

Branch naming conventions supported:
- feature/* → minor version bump
- hotfix/* → patch version bump  
- release/* → prepare for release tagging
- develop/main → analyze commits for bump type`,
	Example: `  esh-cli branch-version --suggest                # Suggest version bump for current branch
  esh-cli branch-version --auto-tag stg6          # Auto-create tag based on branch
  esh-cli branch-version --release-prep           # Prepare for release workflow
  esh-cli branch-version --suggest --service api  # Branch-based suggestion for specific service`,
	Run: runBranchVersion,
}

type BranchInfo struct {
	Name     string
	Type     string
	Feature  string
	Strategy string
}

func init() {
	rootCmd.AddCommand(branchVersionCmd)

	branchVersionCmd.Flags().BoolVar(&branchSuggest, "suggest", false, "Suggest version bump based on branch")
	branchVersionCmd.Flags().BoolVar(&branchAutoTag, "auto-tag", false, "Automatically create appropriate tag")
	branchVersionCmd.Flags().BoolVar(&branchReleasePrep, "release-prep", false, "Prepare release branch workflow")
	branchVersionCmd.Flags().StringVarP(&branchEnvironment, "env", "e", "", "Target environment for tagging")
	branchVersionCmd.Flags().StringVarP(&branchService, "service", "s", "", "Service name for tagging")
}

func runBranchVersion(cmd *cobra.Command, args []string) {
	// Get current branch
	currentBranch, err := utils.Cmd("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
		os.Exit(1)
	}
	currentBranch = strings.TrimSpace(currentBranch)

	branchInfo := analyzeBranch(currentBranch)

	fmt.Printf("🌿 Branch Analysis\n")
	fmt.Printf("Current Branch: %s\n", branchInfo.Name)
	fmt.Printf("Branch Type: %s\n", branchInfo.Type)
	if branchInfo.Feature != "" {
		fmt.Printf("Feature/Issue: %s\n", branchInfo.Feature)
	}
	fmt.Printf("Suggested Strategy: %s\n", branchInfo.Strategy)
	fmt.Println()

	if branchSuggest {
		suggestVersionBump(branchInfo)
	}

	if branchAutoTag {
		if branchEnvironment == "" {
			fmt.Fprintf(os.Stderr, "Error: --env flag is required for auto-tagging\n")
			os.Exit(1)
		}
		autoCreateTag(branchInfo, branchEnvironment, branchService)
	}

	if branchReleasePrep {
		prepareRelease(branchInfo)
	}
}

func analyzeBranch(branchName string) BranchInfo {
	info := BranchInfo{
		Name: branchName,
		Type: "unknown",
	}

	nameLower := strings.ToLower(branchName)

	// Analyze branch type and extract feature name
	if strings.HasPrefix(nameLower, "feature/") {
		info.Type = "feature"
		info.Feature = strings.TrimPrefix(branchName, "feature/")
		info.Strategy = "minor version bump (new features)"
	} else if strings.HasPrefix(nameLower, "feat/") {
		info.Type = "feature"
		info.Feature = strings.TrimPrefix(branchName, "feat/")
		info.Strategy = "minor version bump (new features)"
	} else if strings.HasPrefix(nameLower, "hotfix/") {
		info.Type = "hotfix"
		info.Feature = strings.TrimPrefix(branchName, "hotfix/")
		info.Strategy = "patch version bump (bug fixes)"
	} else if strings.HasPrefix(nameLower, "fix/") {
		info.Type = "hotfix"
		info.Feature = strings.TrimPrefix(branchName, "fix/")
		info.Strategy = "patch version bump (bug fixes)"
	} else if strings.HasPrefix(nameLower, "release/") {
		info.Type = "release"
		info.Feature = strings.TrimPrefix(branchName, "release/")
		info.Strategy = "prepare for release tagging"
	} else if nameLower == "develop" || nameLower == "development" {
		info.Type = "develop"
		info.Strategy = "analyze commits for bump type"
	} else if nameLower == "main" || nameLower == "master" {
		info.Type = "main"
		info.Strategy = "analyze commits for bump type"
	} else if strings.HasPrefix(nameLower, "bugfix/") {
		info.Type = "bugfix"
		info.Feature = strings.TrimPrefix(branchName, "bugfix/")
		info.Strategy = "patch version bump (bug fixes)"
	} else if strings.HasPrefix(nameLower, "chore/") {
		info.Type = "chore"
		info.Feature = strings.TrimPrefix(branchName, "chore/")
		info.Strategy = "patch version bump (maintenance)"
	} else {
		info.Type = "custom"
		info.Strategy = "analyze commits or manual specification"
	}

	return info
}

func suggestVersionBump(branchInfo BranchInfo) {
	fmt.Printf("📋 Version Bump Suggestion\n")

	var bumpType utils.BumpType

	switch branchInfo.Type {
	case "feature":
		bumpType = utils.BumpMinor
		fmt.Printf("Recommended: MINOR bump (new features)\n")
		fmt.Printf("Reason: Feature branch detected\n")
	case "hotfix", "bugfix":
		bumpType = utils.BumpPatch
		fmt.Printf("Recommended: PATCH bump (bug fixes)\n")
		fmt.Printf("Reason: Fix branch detected\n")
	case "release":
		fmt.Printf("Recommended: Review and prepare for release\n")
		fmt.Printf("Reason: Release branch - version should be finalized\n")
		showReleasePreparation(branchInfo)
		return
	case "develop", "main":
		// Analyze commits to suggest bump type
		bumpType = analyzeCommitsForBump()
		fmt.Printf("Recommended: %s bump (based on commit analysis)\n", bumpType)
		fmt.Printf("Reason: Main/develop branch - analyzing recent commits\n")
	case "chore":
		bumpType = utils.BumpPatch
		fmt.Printf("Recommended: PATCH bump (maintenance)\n")
		fmt.Printf("Reason: Chore/maintenance branch\n")
	default:
		bumpType = analyzeCommitsForBump()
		fmt.Printf("Recommended: %s bump (based on commit analysis)\n", bumpType)
		fmt.Printf("Reason: Custom branch - analyzing commits\n")
	}

	// Show example commands
	fmt.Printf("\n💡 Suggested Commands:\n")
	if branchEnvironment != "" {
		fmt.Printf("  esh-cli bump-version %s --%s", branchEnvironment, strings.ToLower(string(bumpType)))
		if branchService != "" {
			fmt.Printf(" --service %s", branchService)
		}
		fmt.Println()
	} else {
		fmt.Printf("  esh-cli bump-version <environment> --%s\n", strings.ToLower(string(bumpType)))
	}

	fmt.Printf("  esh-cli bump-version <environment> --auto  # Auto-detect from commits\n")
	fmt.Printf("  esh-cli bump-version <environment> --preview  # Preview without creating\n")
}

func autoCreateTag(branchInfo BranchInfo, environment, service string) {
	fmt.Printf("🏷️  Auto-Tagging\n")

	// Validate environment
	if !utils.ContainsString(utils.ENVS, environment) {
		fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
			environment, utils.ENVS)
		os.Exit(1)
	}

	var bumpType utils.BumpType

	switch branchInfo.Type {
	case "feature":
		bumpType = utils.BumpMinor
	case "hotfix", "bugfix":
		bumpType = utils.BumpPatch
	case "chore":
		bumpType = utils.BumpPatch
	case "release":
		fmt.Printf("Release branch detected. Use 'esh-cli branch-version --release-prep' instead.\n")
		return
	default:
		bumpType = analyzeCommitsForBump()
	}

	fmt.Printf("Creating %s tag for %s environment...\n", bumpType, environment)

	// Find latest tag
	latestTag, latestVersion, err := utils.GetLatestSemanticVersion(environment, service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding latest version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current latest: %s (%s)\n", latestTag, latestVersion)

	// Create new tag
	newTag, err := utils.BumpTagVersion(latestTag, bumpType, environment, service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new tag: %v\n", err)
		os.Exit(1)
	}

	// Confirm with user
	newVersion, _ := utils.GetVersionFromTag(newTag)
	fmt.Printf("New tag will be: %s (%s)\n", newTag, newVersion)

	if utils.Ask("Create this tag? (y/n)") != "y" {
		fmt.Println("Operation cancelled")
		return
	}

	// Create tag with branch-specific comment
	comment := fmt.Sprintf("Auto-tagged from %s branch: %s", branchInfo.Type, branchInfo.Name)

	// Get current commit
	commit, err := utils.Cmd("git rev-parse HEAD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current commit: %v\n", err)
		os.Exit(1)
	}
	commit = strings.TrimSpace(commit)

	// Create and push tag
	_, err = utils.Cmd(fmt.Sprintf("git tag -a %s -m \"%s\" %s", newTag, comment, commit))
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
}

func prepareRelease(branchInfo BranchInfo) {
	fmt.Printf("🚀 Release Preparation\n")

	if branchInfo.Type == "release" {
		fmt.Printf("Release branch detected: %s\n", branchInfo.Name)

		// Extract version from release branch if possible
		if branchInfo.Feature != "" {
			fmt.Printf("Target version: %s\n", branchInfo.Feature)
		}

		showReleasePreparation(branchInfo)
	} else {
		fmt.Printf("Current branch (%s) is not a release branch.\n", branchInfo.Name)
		fmt.Printf("Consider creating a release branch first:\n")
		fmt.Printf("  git checkout -b release/1.2.0\n")
		fmt.Printf("  git push -u origin release/1.2.0\n")
	}
}

func showReleasePreparation(branchInfo BranchInfo) {
	fmt.Printf("\n📋 Release Checklist:\n")
	fmt.Printf("  ☐ Update version numbers in code\n")
	fmt.Printf("  ☐ Update CHANGELOG.md\n")
	fmt.Printf("  ☐ Run tests: make test\n")
	fmt.Printf("  ☐ Build and verify: make build\n")
	fmt.Printf("  ☐ Create release tags for each environment\n")

	fmt.Printf("\n💡 Suggested Release Commands:\n")
	fmt.Printf("  # Generate changelog\n")
	fmt.Printf("  esh-cli changelog --full --format markdown --output CHANGELOG.md\n")
	fmt.Printf("  \n")
	fmt.Printf("  # Create tags for each environment\n")
	fmt.Printf("  esh-cli add-tag dev <version>\n")
	fmt.Printf("  esh-cli add-tag stg6 <version>\n")
	fmt.Printf("  esh-cli add-tag production2 <version>\n")

	if branchInfo.Feature != "" {
		version := branchInfo.Feature
		fmt.Printf("  \n")
		fmt.Printf("  # For version %s:\n", version)
		fmt.Printf("  esh-cli add-tag dev %s\n", version)
		fmt.Printf("  esh-cli add-tag stg6 %s\n", version)
		fmt.Printf("  esh-cli add-tag production2 %s\n", version)
	}
}

func analyzeCommitsForBump() utils.BumpType {
	// Get recent commits
	output, err := utils.Cmd("git log --oneline -10 --pretty=format:\"%s\"")
	if err != nil {
		return utils.BumpPatch // Default fallback
	}

	if output == "" {
		return utils.BumpPatch
	}

	commits := strings.Split(output, "\n")
	return utils.DetectBumpType(commits)
}
