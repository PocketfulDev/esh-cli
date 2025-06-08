# ESH CLI Semantic Versioning Enhancement Design

## Overview
This document outlines the design for enhanced semantic versioning features in ESH CLI, focusing on git workflow integration and automated version management.

## Current State Analysis
- ✅ Semantic versioning enforcement (MAJOR.MINOR.PATCH)
- ✅ Tag format: `[service_]env_major.minor.patch-release[.hotfix]`
- ✅ Basic increment functionality (release numbers only)
- ✅ Environment promotion capabilities
- ❌ No semantic version bumping (major/minor/patch)
- ❌ No changelog generation
- ❌ No version comparison tools
- ❌ No branch-based versioning strategies

## Proposed New Commands

### 1. `bump-version` Command
**Purpose**: Semantic version bumping with intelligent defaults

**Usage**:
```bash
esh-cli bump-version <environment> [flags]
```

**Flags**:
- `--major`: Bump major version (breaking changes)
- `--minor`: Bump minor version (new features)
- `--patch`: Bump patch version (bug fixes)
- `--auto`: Auto-detect bump type from commit messages
- `--preview`: Show what would be created without executing
- `--service`: Target specific service
- `--from-commit`: Specify commit to tag (default: HEAD)

**Examples**:
```bash
# Manual version bumps
esh-cli bump-version stg6 --major     # 1.2.3 → 2.0.0-1
esh-cli bump-version stg6 --minor     # 1.2.3 → 1.3.0-1  
esh-cli bump-version stg6 --patch     # 1.2.3 → 1.2.4-1

# Auto-detection based on conventional commits
esh-cli bump-version stg6 --auto      # Analyzes commit messages

# Preview mode
esh-cli bump-version stg6 --major --preview
```

### 2. `version-diff` Command  
**Purpose**: Compare versions and analyze semantic differences

**Usage**:
```bash
esh-cli version-diff <tag1> [tag2] [flags]
esh-cli version-diff <environment> [flags]
```

**Flags**:
- `--history`: Show version history
- `--remote`: Compare with remote tags
- `--commits`: Show commits between versions
- `--files`: Show changed files
- `--stats`: Show detailed statistics

### 3. `version-list` Command
**Purpose**: Enhanced tag listing with semantic version awareness

**Usage**:
```bash
esh-cli version-list [environment] [flags]
```

**Flags**:
- `--all`: Show all environments
- `--major <n>`: Filter by major version
- `--minor <n>`: Filter by minor version
- `--format <json|table|compact>`: Output format
- `--sort <version|date>`: Sort order

### 4. `changelog` Command
**Purpose**: Generate changelogs from git history and semantic versions

**Usage**:
```bash
esh-cli changelog <environment> [flags]
esh-cli changelog <tag1>..<tag2> [flags]
```

**Flags**:
- `--format <markdown|json|text>`: Output format
- `--conventional-commits`: Parse conventional commit messages
- `--full`: Generate complete changelog
- `--since <date>`: Changes since date
- `--output <file>`: Write to file

### 5. `branch-version` Command
**Purpose**: Branch-aware versioning and git flow integration

**Usage**:
```bash
esh-cli branch-version [flags]
```

**Flags**:
- `--suggest`: Suggest version bump based on branch
- `--auto-tag`: Automatically create appropriate tag
- `--release-prep`: Prepare release branch workflow

### 6. `version-sync` Command
**Purpose**: Synchronize versions across environments

**Usage**:
```bash
esh-cli version-sync <version> --from <env> --to <env1,env2>
esh-cli version-sync --promote <from-env> <to-env>
```

### 7. `version-validate` Command
**Purpose**: Validation and consistency checks

**Usage**:
```bash
esh-cli version-validate [flags]
```

**Flags**:
- `--pre-release <env> <version>`: Validate before tagging
- `--consistency`: Check cross-environment consistency
- `--fix`: Auto-fix common issues

## Implementation Strategy

### Phase 1: Core Semantic Version Functions
1. **Enhance `pkg/utils/utils.go`**:
   - Add `BumpSemanticVersion()` function
   - Add `ParseSemanticVersion()` function  
   - Add `CompareSemanticVersions()` function
   - Add `ValidateSemanticVersion()` function

2. **Create `pkg/semver/` package**:
   - Dedicated semantic versioning utilities
   - Version comparison and manipulation
   - Changelog parsing utilities

### Phase 2: New Commands Implementation
1. **`cmd/bump-version.go`**: Implement version bumping
2. **`cmd/version-diff.go`**: Version comparison
3. **`cmd/version-list.go`**: Enhanced listing
4. **`cmd/changelog.go`**: Changelog generation

### Phase 3: Advanced Features
1. **Branch integration**: Git flow awareness
2. **Conventional commits**: Commit message parsing
3. **Automation**: CI/CD integration hooks

## Enhanced Utility Functions

### New Functions to Add:

```go
// Semantic version manipulation
func BumpSemanticVersion(version string, bumpType string) string
func ParseSemanticVersion(version string) (major, minor, patch int, err error)
func CompareSemanticVersions(v1, v2 string) int
func IsValidSemanticVersion(version string) bool

// Tag analysis  
func GetVersionFromTag(tag string) string
func GetSemanticVersionFromTag(tag string) (major, minor, patch int, err error)
func FindTagsWithVersion(env, version string) []string

// Git integration
func GetCommitsBetweenTags(tag1, tag2 string) ([]string, error)
func GetBranchType(branch string) string
func SuggestVersionBump(commits []string) string

// Changelog utilities
func ParseConventionalCommit(message string) (type, scope, description string)
func GenerateChangelog(commits []string) string
func GroupCommitsByType(commits []string) map[string][]string
```

## Configuration Enhancements

### New Config Options:
```yaml
# ~/.esh-cli.yaml
versioning:
  conventional_commits: true
  auto_changelog: true
  bump_strategy: "conventional" # conventional, manual, branch-based
  
changelog:
  format: "markdown"
  sections:
    - "Features"
    - "Bug Fixes" 
    - "Breaking Changes"
    - "Documentation"
    
validation:
  require_tests: true
  block_major_on_develop: false
  enforce_semantic_versioning: true

git_flow:
  feature_branch_bump: "minor"
  hotfix_branch_bump: "patch"
  release_branch_bump: "auto"
```

## Backward Compatibility

All new features will be additive and maintain full compatibility with:
- Existing tag formats
- Current command structure  
- Existing validation logic
- Environment configurations

## Testing Strategy

1. **Unit Tests**: For all new utility functions
2. **Integration Tests**: End-to-end command testing
3. **Compatibility Tests**: Ensure existing workflows still work
4. **Performance Tests**: For bulk operations and large repositories

## Benefits

1. **Enhanced Developer Experience**: Intuitive semantic versioning
2. **Automation**: Reduce manual version management overhead
3. **Consistency**: Enforce semantic versioning across all environments
4. **Visibility**: Better version tracking and change documentation
5. **Git Flow Integration**: Support modern development workflows
6. **CI/CD Ready**: Features designed for automated pipelines

## Future Considerations

1. **Plugin System**: Allow custom version bump strategies
2. **Integration**: GitHub/GitLab release integration
3. **Notifications**: Slack/Teams integration for version updates
4. **Analytics**: Version deployment success tracking
5. **Multi-repo**: Support for monorepo version management
