package cmd

import (
	"esh-cli/pkg/utils"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	changelogFormat          string
	changelogConventional    bool
	changelogFull            bool
	changelogSince           string
	changelogOutput          string
	changelogFromTag         string
	changelogToTag           string
	changelogGroupByType     bool
	changelogIncludeBreaking bool
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog [environment] [flags]",
	Short: "Generate changelogs from git history and semantic versions",
	Long: `Generate changelogs from git history and semantic versions.

This command analyzes git commits and creates formatted changelogs:
- Supports conventional commit parsing
- Groups changes by type (features, fixes, breaking changes)
- Generates markdown, JSON, or text output
- Supports date ranges and tag ranges

The changelog can be generated for a specific environment, between two tags,
or since a specific date.`,
	Example: `  esh-cli changelog stg6                           # Generate changelog for staging
  esh-cli changelog --from stg6_1.2.0-1 --to stg6_1.3.0-1  # Between specific tags
  esh-cli changelog --since 2024-01-01            # Changes since date
  esh-cli changelog stg6 --conventional-commits   # Parse conventional commits
  esh-cli changelog stg6 --format json --output changelog.json`,
	Args: cobra.MaximumNArgs(1),
	Run:  runChangelog,
}

type ChangelogEntry struct {
	Type        string
	Scope       string
	Description string
	Hash        string
	Breaking    bool
	Date        time.Time
}

type Changelog struct {
	Title     string
	FromTag   string
	ToTag     string
	FromDate  time.Time
	ToDate    time.Time
	Entries   []ChangelogEntry
	GroupedBy map[string][]ChangelogEntry
}

func init() {
	rootCmd.AddCommand(changelogCmd)

	changelogCmd.Flags().StringVar(&changelogFormat, "format", "markdown", "Output format (markdown, json, text)")
	changelogCmd.Flags().BoolVar(&changelogConventional, "conventional-commits", false, "Parse conventional commit messages")
	changelogCmd.Flags().BoolVar(&changelogFull, "full", false, "Generate complete changelog")
	changelogCmd.Flags().StringVar(&changelogSince, "since", "", "Changes since date (YYYY-MM-DD)")
	changelogCmd.Flags().StringVar(&changelogOutput, "output", "", "Write to file (default: stdout)")
	changelogCmd.Flags().StringVar(&changelogFromTag, "from", "", "Start tag for range")
	changelogCmd.Flags().StringVar(&changelogToTag, "to", "", "End tag for range")
	changelogCmd.Flags().BoolVar(&changelogGroupByType, "group-by-type", true, "Group entries by type")
	changelogCmd.Flags().BoolVar(&changelogIncludeBreaking, "include-breaking", true, "Include breaking changes section")
}

func runChangelog(cmd *cobra.Command, args []string) {
	var environment string
	if len(args) > 0 {
		environment = args[0]
		if !utils.ContainsString(utils.ENVS, environment) {
			fmt.Fprintf(os.Stderr, "Error: invalid environment '%s'. Valid environments: %v\n",
				environment, utils.ENVS)
			os.Exit(1)
		}
	}

	// Determine tag range
	fromTag := changelogFromTag
	toTag := changelogToTag

	if environment != "" && fromTag == "" && toTag == "" {
		// Get latest and previous tag for environment
		latest, err := getLatestTagForEnvironment(environment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding latest tag: %v\n", err)
			os.Exit(1)
		}
		toTag = latest

		if !changelogFull {
			previous, err := findPreviousTag(latest)
			if err == nil {
				fromTag = previous
			}
		}
	}

	// Generate changelog
	changelog, err := generateChangelog(fromTag, toTag, environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating changelog: %v\n", err)
		os.Exit(1)
	}

	// Format output
	var output string
	switch changelogFormat {
	case "markdown":
		output = formatMarkdown(changelog)
	case "json":
		output = formatJSON(changelog)
	case "text":
		output = formatText(changelog)
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported format '%s'. Use: markdown, json, text\n", changelogFormat)
		os.Exit(1)
	}

	// Write output
	if changelogOutput != "" {
		err := os.WriteFile(changelogOutput, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Changelog written to: %s\n", changelogOutput)
	} else {
		fmt.Print(output)
	}
}

func generateChangelog(fromTag, toTag, environment string) (*Changelog, error) {
	changelog := &Changelog{
		FromTag:   fromTag,
		ToTag:     toTag,
		GroupedBy: make(map[string][]ChangelogEntry),
	}

	// Set title
	if environment != "" {
		if fromTag != "" && toTag != "" {
			changelog.Title = fmt.Sprintf("Changelog for %s (%s → %s)", environment, fromTag, toTag)
		} else if toTag != "" {
			changelog.Title = fmt.Sprintf("Changelog for %s (%s)", environment, toTag)
		} else {
			changelog.Title = fmt.Sprintf("Changelog for %s", environment)
		}
	} else if fromTag != "" && toTag != "" {
		changelog.Title = fmt.Sprintf("Changelog (%s → %s)", fromTag, toTag)
	} else {
		changelog.Title = "Changelog"
	}

	// Get commits
	var commits []string
	var err error

	if fromTag != "" && toTag != "" {
		commits, err = utils.GetCommitsBetweenTags(fromTag, toTag)
	} else if toTag != "" {
		// Get all commits up to toTag
		output, e := utils.Cmd(fmt.Sprintf("git log --oneline --pretty=format:\"%%H %%s\" %s", toTag))
		if e != nil {
			err = e
		} else if output != "" {
			commits = strings.Split(output, "\n")
		}
	} else if changelogSince != "" {
		// Get commits since date
		output, e := utils.Cmd(fmt.Sprintf("git log --oneline --pretty=format:\"%%H %%s\" --since=\"%s\"", changelogSince))
		if e != nil {
			err = e
		} else if output != "" {
			commits = strings.Split(output, "\n")
		}
	} else {
		// Get recent commits
		output, e := utils.Cmd("git log --oneline --pretty=format:\"%H %s\" -20")
		if e != nil {
			err = e
		} else if output != "" {
			commits = strings.Split(output, "\n")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error getting commits: %v", err)
	}

	// Parse commits into changelog entries
	for _, commit := range commits {
		if commit == "" {
			continue
		}

		entry := parseCommit(commit)
		if entry != nil {
			changelog.Entries = append(changelog.Entries, *entry)

			// Group by type
			if changelogGroupByType {
				typeKey := entry.Type
				if typeKey == "" {
					typeKey = "Other"
				}
				changelog.GroupedBy[typeKey] = append(changelog.GroupedBy[typeKey], *entry)
			}
		}
	}

	// Sort entries by date (newest first)
	sort.Slice(changelog.Entries, func(i, j int) bool {
		return changelog.Entries[i].Date.After(changelog.Entries[j].Date)
	})

	return changelog, nil
}

func parseCommit(commit string) *ChangelogEntry {
	parts := strings.SplitN(commit, " ", 2)
	if len(parts) != 2 {
		return nil
	}

	hash := parts[0]
	message := parts[1]

	entry := &ChangelogEntry{
		Hash:        hash,
		Description: message,
		Date:        getCommitDate(hash),
	}

	if changelogConventional {
		parseConventionalCommit(entry, message)
	} else {
		// Simple parsing - try to detect type from message
		entry.Type = detectCommitType(message)
	}

	return entry
}

func parseConventionalCommit(entry *ChangelogEntry, message string) {
	// Conventional commit format: type(scope): description
	// Optional: type(scope)!: description (breaking change)
	conventionalRegex := regexp.MustCompile(`^(\w+)(\([^)]+\))?(!)?: (.+)$`)
	matches := conventionalRegex.FindStringSubmatch(message)

	if len(matches) >= 5 {
		entry.Type = matches[1]
		if matches[2] != "" {
			// Remove parentheses from scope
			entry.Scope = strings.Trim(matches[2], "()")
		}
		entry.Breaking = matches[3] == "!"
		entry.Description = matches[4]
	} else {
		// Fallback to simple type detection
		entry.Type = detectCommitType(message)
		entry.Description = message
	}

	// Check for BREAKING CHANGE in message body
	if strings.Contains(strings.ToUpper(message), "BREAKING CHANGE") {
		entry.Breaking = true
	}
}

func detectCommitType(message string) string {
	messageLower := strings.ToLower(message)

	if strings.HasPrefix(messageLower, "feat") || strings.Contains(messageLower, "feature") || strings.Contains(messageLower, "add ") {
		return "feat"
	}
	if strings.HasPrefix(messageLower, "fix") || strings.Contains(messageLower, "bug") || strings.Contains(messageLower, "issue") {
		return "fix"
	}
	if strings.Contains(messageLower, "doc") || strings.Contains(messageLower, "readme") {
		return "docs"
	}
	if strings.Contains(messageLower, "style") || strings.Contains(messageLower, "format") {
		return "style"
	}
	if strings.Contains(messageLower, "refactor") || strings.Contains(messageLower, "cleanup") {
		return "refactor"
	}
	if strings.Contains(messageLower, "test") {
		return "test"
	}
	if strings.Contains(messageLower, "chore") || strings.Contains(messageLower, "update") {
		return "chore"
	}

	return "other"
}

func getCommitDate(hash string) time.Time {
	output, err := utils.Cmd(fmt.Sprintf("git show -s --format=%%ci %s", hash))
	if err != nil {
		return time.Time{}
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(output))
	if err != nil {
		return time.Time{}
	}

	return date
}

func formatMarkdown(changelog *Changelog) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", changelog.Title))

	if changelog.FromTag != "" && changelog.ToTag != "" {
		sb.WriteString(fmt.Sprintf("**Full Changelog**: %s...%s\n\n", changelog.FromTag, changelog.ToTag))
	}

	if changelogGroupByType && len(changelog.GroupedBy) > 0 {
		// Define order for sections
		sectionOrder := []string{"feat", "fix", "perf", "refactor", "docs", "style", "test", "chore", "other"}
		sectionTitles := map[string]string{
			"feat":     "🚀 Features",
			"fix":      "🐛 Bug Fixes",
			"perf":     "⚡ Performance",
			"refactor": "♻️ Refactoring",
			"docs":     "📚 Documentation",
			"style":    "💄 Style",
			"test":     "🧪 Tests",
			"chore":    "🔧 Chores",
			"other":    "📝 Other Changes",
		}

		// Add breaking changes first if any
		if changelogIncludeBreaking {
			breakingChanges := getBreakingChanges(changelog.Entries)
			if len(breakingChanges) > 0 {
				sb.WriteString("## 💥 Breaking Changes\n\n")
				for _, entry := range breakingChanges {
					sb.WriteString(formatMarkdownEntry(entry))
				}
				sb.WriteString("\n")
			}
		}

		// Add sections in order
		for _, sectionKey := range sectionOrder {
			if entries, exists := changelog.GroupedBy[sectionKey]; exists && len(entries) > 0 {
				title := sectionTitles[sectionKey]
				if title == "" {
					title = strings.Title(sectionKey)
				}

				sb.WriteString(fmt.Sprintf("## %s\n\n", title))
				for _, entry := range entries {
					if !entry.Breaking || !changelogIncludeBreaking {
						sb.WriteString(formatMarkdownEntry(entry))
					}
				}
				sb.WriteString("\n")
			}
		}
	} else {
		// Simple list format
		sb.WriteString("## Changes\n\n")
		for _, entry := range changelog.Entries {
			sb.WriteString(formatMarkdownEntry(entry))
		}
	}

	return sb.String()
}

func formatMarkdownEntry(entry ChangelogEntry) string {
	var sb strings.Builder

	sb.WriteString("- ")

	if entry.Scope != "" {
		sb.WriteString(fmt.Sprintf("**%s**: ", entry.Scope))
	}

	sb.WriteString(entry.Description)

	if entry.Breaking {
		sb.WriteString(" ⚠️ **BREAKING**")
	}

	sb.WriteString(fmt.Sprintf(" ([%s])\n", entry.Hash[:8]))

	return sb.String()
}

func formatJSON(changelog *Changelog) string {
	// Simple JSON format - in a real implementation you'd use json.Marshal
	var sb strings.Builder

	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"title\": \"%s\",\n", changelog.Title))
	sb.WriteString(fmt.Sprintf("  \"from_tag\": \"%s\",\n", changelog.FromTag))
	sb.WriteString(fmt.Sprintf("  \"to_tag\": \"%s\",\n", changelog.ToTag))
	sb.WriteString("  \"entries\": [\n")

	for i, entry := range changelog.Entries {
		sb.WriteString("    {\n")
		sb.WriteString(fmt.Sprintf("      \"type\": \"%s\",\n", entry.Type))
		sb.WriteString(fmt.Sprintf("      \"scope\": \"%s\",\n", entry.Scope))
		sb.WriteString(fmt.Sprintf("      \"description\": \"%s\",\n", entry.Description))
		sb.WriteString(fmt.Sprintf("      \"hash\": \"%s\",\n", entry.Hash))
		sb.WriteString(fmt.Sprintf("      \"breaking\": %t,\n", entry.Breaking))
		sb.WriteString(fmt.Sprintf("      \"date\": \"%s\"\n", entry.Date.Format(time.RFC3339)))
		if i < len(changelog.Entries)-1 {
			sb.WriteString("    },\n")
		} else {
			sb.WriteString("    }\n")
		}
	}

	sb.WriteString("  ]\n")
	sb.WriteString("}\n")

	return sb.String()
}

func formatText(changelog *Changelog) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", changelog.Title))
	sb.WriteString(strings.Repeat("=", len(changelog.Title)) + "\n\n")

	if changelog.FromTag != "" && changelog.ToTag != "" {
		sb.WriteString(fmt.Sprintf("Range: %s → %s\n\n", changelog.FromTag, changelog.ToTag))
	}

	for _, entry := range changelog.Entries {
		sb.WriteString(fmt.Sprintf("* %s", entry.Description))
		if entry.Scope != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", entry.Scope))
		}
		if entry.Breaking {
			sb.WriteString(" [BREAKING]")
		}
		sb.WriteString(fmt.Sprintf(" [%s]\n", entry.Hash[:8]))
	}

	return sb.String()
}

func getLatestTagForEnvironment(environment string) (string, error) {
	pattern := fmt.Sprintf("*_%s_*", environment)
	output, err := utils.Cmd(fmt.Sprintf("git tag -l '%s' --sort=-version:refname", pattern))
	if err != nil {
		return "", err
	}

	if output == "" {
		return "", fmt.Errorf("no tags found for environment: %s", environment)
	}

	tags := strings.Split(output, "\n")
	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found for environment: %s", environment)
	}

	return tags[0], nil
}

func getBreakingChanges(entries []ChangelogEntry) []ChangelogEntry {
	var breaking []ChangelogEntry
	for _, entry := range entries {
		if entry.Breaking {
			breaking = append(breaking, entry)
		}
	}
	return breaking
}
