package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// BumpType represents the type of semantic version bump
type BumpType string

const (
	BumpMajor BumpType = "major"
	BumpMinor BumpType = "minor"
	BumpPatch BumpType = "patch"
	BumpAuto  BumpType = "auto"
)

// SemanticVersion represents a parsed semantic version
type SemanticVersion struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

// String returns the string representation of the semantic version
func (sv SemanticVersion) String() string {
	version := fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
	if sv.Prerelease != "" {
		version += "-" + sv.Prerelease
	}
	return version
}

// ParseSemanticVersion parses a semantic version string into its components
func ParseSemanticVersion(version string) (*SemanticVersion, error) {
	// Remove any prefix (like "v")
	version = strings.TrimPrefix(version, "v")

	// Split by prerelease separator
	parts := strings.Split(version, "-")
	versionPart := parts[0]
	prerelease := ""
	if len(parts) > 1 {
		prerelease = strings.Join(parts[1:], "-")
	}

	// Parse version numbers
	versionNumbers := strings.Split(versionPart, ".")
	if len(versionNumbers) != 3 {
		return nil, fmt.Errorf("invalid semantic version format: %s", version)
	}

	major, err := strconv.Atoi(versionNumbers[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", versionNumbers[0])
	}

	minor, err := strconv.Atoi(versionNumbers[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", versionNumbers[1])
	}

	patch, err := strconv.Atoi(versionNumbers[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", versionNumbers[2])
	}

	return &SemanticVersion{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
	}, nil
}

// BumpSemanticVersion bumps a semantic version according to the specified bump type
func BumpSemanticVersion(version string, bumpType BumpType) (string, error) {
	sv, err := ParseSemanticVersion(version)
	if err != nil {
		return "", err
	}

	switch bumpType {
	case BumpMajor:
		sv.Major++
		sv.Minor = 0
		sv.Patch = 0
	case BumpMinor:
		sv.Minor++
		sv.Patch = 0
	case BumpPatch:
		sv.Patch++
	default:
		return "", fmt.Errorf("unsupported bump type: %s", bumpType)
	}

	// Remove prerelease for bumped versions
	sv.Prerelease = ""

	return sv.String(), nil
}

// CompareSemanticVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareSemanticVersions(v1, v2 string) (int, error) {
	sv1, err := ParseSemanticVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("error parsing version 1: %v", err)
	}

	sv2, err := ParseSemanticVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("error parsing version 2: %v", err)
	}

	// Compare major version
	if sv1.Major != sv2.Major {
		if sv1.Major < sv2.Major {
			return -1, nil
		}
		return 1, nil
	}

	// Compare minor version
	if sv1.Minor != sv2.Minor {
		if sv1.Minor < sv2.Minor {
			return -1, nil
		}
		return 1, nil
	}

	// Compare patch version
	if sv1.Patch != sv2.Patch {
		if sv1.Patch < sv2.Patch {
			return -1, nil
		}
		return 1, nil
	}

	// Versions are equal (ignoring prerelease for now)
	return 0, nil
}

// GetVersionFromTag extracts the semantic version from a tag
func GetVersionFromTag(tag string) (string, error) {
	if !IsTagValid(tag) {
		return "", fmt.Errorf("invalid tag format: %s", tag)
	}

	parts := strings.Split(tag, "_")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid tag format: %s", tag)
	}

	versionPart := parts[len(parts)-1]
	versionReleaseParts := strings.Split(versionPart, "-")

	return versionReleaseParts[0], nil
}

// BumpTagVersion creates a new tag with bumped semantic version
func BumpTagVersion(tag string, bumpType BumpType, env string, service string) (string, error) {
	version, err := GetVersionFromTag(tag)
	if err != nil {
		return "", err
	}

	newVersion, err := BumpSemanticVersion(version, bumpType)
	if err != nil {
		return "", err
	}

	// Create new tag with bumped version and release suffix -1
	newTag := fmt.Sprintf("%s-1", TagPrefix(env, newVersion, service))
	return newTag, nil
}

// GetCommitsBetweenTags gets commit messages between two tags
func GetCommitsBetweenTags(tag1, tag2 string) ([]string, error) {
	if tag1 == "" || tag2 == "" {
		return nil, fmt.Errorf("both tags must be provided")
	}

	// Use git log to get commits between tags
	output, err := Cmd(fmt.Sprintf("git log --oneline --pretty=format:\"%%s\" %s..%s", tag1, tag2))
	if err != nil {
		return nil, fmt.Errorf("error getting commits between tags: %v", err)
	}

	if output == "" {
		return []string{}, nil
	}

	commits := strings.Split(output, "\n")
	return commits, nil
}

// DetectBumpType analyzes commit messages to suggest version bump type
func DetectBumpType(commits []string) BumpType {
	hasBreaking := false
	hasFeature := false
	hasFix := false

	// Conventional commit patterns
	breakingPattern := regexp.MustCompile(`(?i)(BREAKING|BREAKING CHANGE|!:)`)
	featurePattern := regexp.MustCompile(`(?i)^(feat|feature)(\(.+\))?:`)
	fixPattern := regexp.MustCompile(`(?i)^(fix|bugfix)(\(.+\))?:`)

	for _, commit := range commits {
		if breakingPattern.MatchString(commit) {
			hasBreaking = true
		} else if featurePattern.MatchString(commit) {
			hasFeature = true
		} else if fixPattern.MatchString(commit) {
			hasFix = true
		}
	}

	// Determine bump type based on commit analysis
	if hasBreaking {
		return BumpMajor
	} else if hasFeature {
		return BumpMinor
	} else if hasFix {
		return BumpPatch
	}

	// Default to patch for any other changes
	return BumpPatch
}

// ValidateSemanticVersionBump validates if a version bump is appropriate
func ValidateSemanticVersionBump(currentTag, proposedVersion string, bumpType BumpType) error {
	currentVersion, err := GetVersionFromTag(currentTag)
	if err != nil {
		return fmt.Errorf("error parsing current tag: %v", err)
	}

	expectedVersion, err := BumpSemanticVersion(currentVersion, bumpType)
	if err != nil {
		return fmt.Errorf("error calculating expected version: %v", err)
	}

	if proposedVersion != expectedVersion {
		return fmt.Errorf("proposed version %s does not match expected %s for %s bump",
			proposedVersion, expectedVersion, bumpType)
	}

	return nil
}

// GetLatestSemanticVersion finds the latest semantic version for an environment
func GetLatestSemanticVersion(env string, service string) (string, string, error) {
	// Get all tags for the environment
	pattern := env
	if service != "" {
		pattern = service + "_" + env
	}

	output, err := Cmd(fmt.Sprintf("git tag -l '%s_*' --sort=-version:refname", pattern))
	if err != nil {
		return "", "", fmt.Errorf("error listing tags: %v", err)
	}

	if output == "" {
		return "", "", fmt.Errorf("no tags found for environment: %s", env)
	}

	tags := strings.Split(output, "\n")
	if len(tags) == 0 {
		return "", "", fmt.Errorf("no tags found for environment: %s", env)
	}

	// Return the first (latest) tag and its version
	latestTag := tags[0]
	version, err := GetVersionFromTag(latestTag)
	if err != nil {
		return "", "", fmt.Errorf("error parsing latest tag: %v", err)
	}

	return latestTag, version, nil
}
