package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	listAll    bool
	listMajor  int
	listMinor  int
	listFormat string
	listSort   string
	listLimit  int
)

// versionListCmd represents the version-list command
var versionListCmd = &cobra.Command{
	Use:   "version-list [environment]",
	Short: "List semantic versions with enhanced filtering",
	Long: `List semantic versions for environments with advanced filtering and formatting options.

This command provides enhanced listing capabilities for tags with semantic version awareness,
including filtering by version components, cross-environment comparison, and multiple output formats.`,
	Example: `  esh-cli version-list stg6                    # List all versions for stg6
  esh-cli version-list stg6 --major 1         # Filter by major version 1
  esh-cli version-list stg6 --major 1 --minor 2  # Filter by version 1.2.x
  esh-cli version-list --all                  # Compare all environments
  esh-cli version-list stg6 --format json     # Output as JSON
  esh-cli version-list stg6 --sort date       # Sort by date instead of version`,
	Args: cobra.MaximumNArgs(1),
	Run:  runVersionList,
}

func init() {
	rootCmd.AddCommand(versionListCmd)

	versionListCmd.Flags().BoolVar(&listAll, "all", false, "show all environments")
	versionListCmd.Flags().IntVar(&listMajor, "major", -1, "filter by major version")
	versionListCmd.Flags().IntVar(&listMinor, "minor", -1, "filter by minor version")
	versionListCmd.Flags().StringVar(&listFormat, "format", "table", "output format (table, json, compact)")
	versionListCmd.Flags().StringVar(&listSort, "sort", "version", "sort order (version, date)")
	versionListCmd.Flags().IntVar(&listLimit, "limit", 10, "maximum number of results per environment")
}

type VersionInfo struct {
	Tag         string    `json:"tag"`
	Environment string    `json:"environment"`
	Service     string    `json:"service"`
	Version     string    `json:"version"`
	Major       int       `json:"major"`
	Minor       int       `json:"minor"`
	Patch       int       `json:"patch"`
	Release     string    `json:"release"`
	Date        time.Time `json:"date"`
	Commit      string    `json:"commit"`
	Message     string    `json:"message"`
}

func runVersionList(cmd *cobra.Command, args []string) {
	var environments []string

	if listAll {
		environments = utils.ENVS
	} else {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: must specify environment or use --all flag\n")
			os.Exit(1)
		}
		environment := args[0]
		if !utils.ContainsString(utils.ENVS, environment) {
			fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
				environment, utils.ENVS)
			os.Exit(1)
		}
		environments = []string{environment}
	}

	// Collect version information
	var allVersions []VersionInfo

	for _, env := range environments {
		versions, err := getVersionsForEnvironment(env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error getting versions for %s: %v\n", env, err)
			continue
		}
		allVersions = append(allVersions, versions...)
	}

	// Apply filters
	filteredVersions := applyFilters(allVersions)

	// Sort results
	sortVersions(filteredVersions)

	// Apply limit
	if listLimit > 0 && len(filteredVersions) > listLimit {
		if !listAll {
			filteredVersions = filteredVersions[:listLimit]
		}
	}

	// Output results
	switch listFormat {
	case "json":
		outputJSON(filteredVersions)
	case "compact":
		outputCompact(filteredVersions)
	default:
		outputTable(filteredVersions)
	}
}

func getVersionsForEnvironment(env string) ([]VersionInfo, error) {
	// Get all tags for environment
	pattern := fmt.Sprintf("*_%s_*", env)
	output, err := utils.Cmd(fmt.Sprintf("git tag -l '%s' --sort=-version:refname", pattern))
	if err != nil {
		return nil, fmt.Errorf("error listing tags: %v", err)
	}

	if output == "" {
		return []VersionInfo{}, nil
	}

	tags := strings.Split(output, "\n")
	var versions []VersionInfo

	for _, tag := range tags {
		if tag == "" {
			continue
		}

		versionInfo, err := parseVersionInfo(tag, env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid tag %s: %v\n", tag, err)
			continue
		}

		versions = append(versions, versionInfo)
	}

	return versions, nil
}

func parseVersionInfo(tag, env string) (VersionInfo, error) {
	if !utils.IsTagValid(tag) {
		return VersionInfo{}, fmt.Errorf("invalid tag format")
	}

	// Parse tag components
	parts := strings.Split(tag, "_")
	var service string
	var envPart string
	var versionPart string

	if len(parts) == 3 {
		// service_env_version format
		service = parts[0]
		envPart = parts[1]
		versionPart = parts[2]
	} else if len(parts) == 2 {
		// env_version format
		envPart = parts[0]
		versionPart = parts[1]
	} else {
		return VersionInfo{}, fmt.Errorf("unexpected tag format")
	}

	if envPart != env {
		return VersionInfo{}, fmt.Errorf("environment mismatch")
	}

	// Parse version and release
	versionReleaseParts := strings.Split(versionPart, "-")
	version := versionReleaseParts[0]
	release := ""
	if len(versionReleaseParts) > 1 {
		release = strings.Join(versionReleaseParts[1:], "-")
	}

	// Parse semantic version
	sv, err := utils.ParseSemanticVersion(version)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("error parsing semantic version: %v", err)
	}

	// Get tag metadata
	commit, err := utils.Cmd(fmt.Sprintf("git rev-list -n 1 %s", tag))
	if err != nil {
		commit = "unknown"
	}

	// Get tag date
	dateStr, err := utils.Cmd(fmt.Sprintf("git log -1 --format=%%ai %s", tag))
	if err != nil {
		dateStr = ""
	}

	var tagDate time.Time
	if dateStr != "" {
		tagDate, _ = time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	}

	// Get tag message
	message, err := utils.Cmd(fmt.Sprintf("git tag -l --format='%%(contents)' %s", tag))
	if err != nil {
		message = ""
	}

	return VersionInfo{
		Tag:         tag,
		Environment: env,
		Service:     service,
		Version:     version,
		Major:       sv.Major,
		Minor:       sv.Minor,
		Patch:       sv.Patch,
		Release:     release,
		Date:        tagDate,
		Commit:      commit,
		Message:     strings.TrimSpace(message),
	}, nil
}

func applyFilters(versions []VersionInfo) []VersionInfo {
	var filtered []VersionInfo

	for _, v := range versions {
		// Apply major version filter
		if listMajor >= 0 && v.Major != listMajor {
			continue
		}

		// Apply minor version filter
		if listMinor >= 0 && v.Minor != listMinor {
			continue
		}

		filtered = append(filtered, v)
	}

	return filtered
}

func sortVersions(versions []VersionInfo) {
	if listSort == "date" {
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].Date.After(versions[j].Date)
		})
	} else {
		// Sort by semantic version (descending)
		sort.Slice(versions, func(i, j int) bool {
			vi := versions[i]
			vj := versions[j]

			if vi.Major != vj.Major {
				return vi.Major > vj.Major
			}
			if vi.Minor != vj.Minor {
				return vi.Minor > vj.Minor
			}
			if vi.Patch != vj.Patch {
				return vi.Patch > vj.Patch
			}

			// For same semantic version, compare release numbers
			releaseI, _ := strconv.Atoi(vi.Release)
			releaseJ, _ := strconv.Atoi(vj.Release)
			return releaseI > releaseJ
		})
	}
}

func outputTable(versions []VersionInfo) {
	if len(versions) == 0 {
		fmt.Println("No versions found matching criteria")
		return
	}

	// Group by environment for better display
	envGroups := make(map[string][]VersionInfo)
	for _, v := range versions {
		envGroups[v.Environment] = append(envGroups[v.Environment], v)
	}

	for env, envVersions := range envGroups {
		fmt.Printf("\n🏷️  Environment: %s\n", env)
		fmt.Printf("%-25s %-12s %-8s %-20s %-12s\n", "Tag", "Version", "Release", "Date", "Commit")
		fmt.Println(strings.Repeat("-", 85))

		for _, v := range envVersions {
			dateStr := v.Date.Format("2006-01-02 15:04")
			if v.Date.IsZero() {
				dateStr = "unknown"
			}

			commitShort := v.Commit
			if len(commitShort) > 8 {
				commitShort = commitShort[:8]
			}

			fmt.Printf("%-25s %-12s %-8s %-20s %-12s\n",
				v.Tag, v.Version, v.Release, dateStr, commitShort)
		}
	}
}

func outputCompact(versions []VersionInfo) {
	for _, v := range versions {
		fmt.Printf("%s (%s)\n", v.Tag, v.Version)
	}
}

func outputJSON(versions []VersionInfo) {
	// Simple JSON output (in a real implementation, use proper JSON marshaling)
	fmt.Println("[")
	for i, v := range versions {
		fmt.Printf("  {\n")
		fmt.Printf("    \"tag\": \"%s\",\n", v.Tag)
		fmt.Printf("    \"environment\": \"%s\",\n", v.Environment)
		fmt.Printf("    \"service\": \"%s\",\n", v.Service)
		fmt.Printf("    \"version\": \"%s\",\n", v.Version)
		fmt.Printf("    \"major\": %d,\n", v.Major)
		fmt.Printf("    \"minor\": %d,\n", v.Minor)
		fmt.Printf("    \"patch\": %d,\n", v.Patch)
		fmt.Printf("    \"release\": \"%s\",\n", v.Release)
		fmt.Printf("    \"commit\": \"%s\",\n", v.Commit)
		fmt.Printf("    \"message\": \"%s\"\n", v.Message)
		if i < len(versions)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")
}
