package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Environment constants
var ENVS = []string{"dev", "mimic2", "stg6", "demo", "production2"}

// Regex patterns
var (
	VersionHotFixPattern     = regexp.MustCompile(`^(\d+)\.(\d+)-(\d+)\.(\d+)`)
	VersionPattern           = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	ReleasePattern           = regexp.MustCompile(`^(\d+)$`)
	ReleasePatternWithHotFix = regexp.MustCompile(`^(\d+).(\d+)$`)
	ReleaseBranchPattern     = regexp.MustCompile(`^release_(\d+)\.(\d+)`)
)

// Cmd executes a shell command and returns the trimmed output
func Cmd(command string) (string, error) {
	fmt.Printf("> %s\n", command)
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Ask prompts the user for input
func Ask(prompt string) string {
	fmt.Printf("\n%s :", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input))
}

// IsVersionValid checks if version string matches the expected pattern
func IsVersionValid(version string, hotFix bool) bool {
	if hotFix {
		return VersionHotFixPattern.MatchString(version)
	}
	return VersionPattern.MatchString(version)
}

// IsTagValid validates a tag format (env_version-release)
func IsTagValid(tag string) bool {
	parts := strings.Split(tag, "_")
	if len(parts) < 2 {
		return false
	}

	env := parts[0]
	if len(parts) == 3 {
		env = parts[1] // service_env_version format
	}

	// Check if env is valid
	envValid := false
	for _, validEnv := range ENVS {
		if env == validEnv {
			envValid = true
			break
		}
	}
	if !envValid {
		return false
	}

	// Get version part
	versionPart := parts[len(parts)-1]
	versionReleaseParts := strings.Split(versionPart, "-")
	if len(versionReleaseParts) != 2 {
		return false
	}

	version := versionReleaseParts[0]
	release := versionReleaseParts[1]

	return VersionPattern.MatchString(version) &&
		(ReleasePattern.MatchString(release) || ReleasePatternWithHotFix.MatchString(release))
}

// GetToday returns today's date in YYYYMMDD format
func GetToday() string {
	return time.Now().Format("20060102")
}

// IncrementTag increments a tag version
func IncrementTag(tag string, hotFix bool) string {
	if tag == "" {
		return ""
	}

	// Validate tag format first
	if !IsTagValid(tag) {
		return ""
	}

	parts := strings.Split(tag, "_")
	if len(parts) < 2 {
		return ""
	}

	versionPart := parts[len(parts)-1]
	prefix := strings.Join(parts[:len(parts)-1], "_")

	versionReleaseParts := strings.Split(versionPart, "-")
	if len(versionReleaseParts) != 2 {
		return ""
	}

	version := versionReleaseParts[0]
	release := versionReleaseParts[1]

	prerelease := "0"

	// Detect tag with prerelease value (hot fix)
	if strings.Contains(release, ".") {
		releaseParts := strings.Split(release, ".")
		if len(releaseParts) != 2 {
			return ""
		}
		release = releaseParts[0]
		prerelease = releaseParts[1]
	}

	if !hotFix {
		releaseInt, err := strconv.Atoi(release)
		if err != nil {
			return ""
		}
		release = strconv.Itoa(releaseInt + 1)
		return fmt.Sprintf("%s_%s-%s", prefix, version, release)
	} else {
		prereleaseInt, err := strconv.Atoi(prerelease)
		if err != nil {
			return ""
		}
		prerelease = strconv.Itoa(prereleaseInt + 1)
		return fmt.Sprintf("%s_%s-%s.%s", prefix, version, release, prerelease)
	}
}

// TagPrefix generates a tag prefix
func TagPrefix(env, version, service string) string {
	prefix := fmt.Sprintf("%s_%s", env, version)
	if service != "" {
		prefix = fmt.Sprintf("%s_%s", service, prefix)
	}
	return prefix
}

// GetEnvFromTag extracts environment from tag
func GetEnvFromTag(tag string) (string, error) {
	parts := strings.Split(tag, "_")
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("tag must have 2 or 3 parts")
	}

	if len(parts) == 2 {
		return parts[0], nil
	}
	return parts[1], nil
}

// FindLastTagAndComment finds the last tag and its comment
func FindLastTagAndComment(env, version, service string) (string, string, error) {
	tagPattern := TagPrefix(env, version, service) + "*"
	command := fmt.Sprintf("git tag --list -n1 '%s' --sort=creatordate | tail -1", tagPattern)

	tagComment, err := Cmd(command)
	if err != nil || tagComment == "" {
		return "", "", err
	}

	parts := strings.SplitN(tagComment, " ", 2)
	tag := parts[0]
	comment := ""
	if len(parts) > 1 {
		comment = strings.TrimSpace(parts[1])
	}

	// Validate that the tag is in the correct format before returning it
	if !IsTagValid(tag) {
		return "", "", nil // Return empty if tag format is invalid
	}

	return tag, comment, nil
}

// IsReleaseBranch checks if branch is a release branch
func IsReleaseBranch(branch string) bool {
	return ReleaseBranchPattern.MatchString(branch)
}

// ContainsString checks if a slice contains a string
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
