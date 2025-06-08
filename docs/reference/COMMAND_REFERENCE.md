# ESH CLI Command Reference

## 🚀 Semantic Versioning Commands

### `bump-version` - Intelligent Version Bumping

**Purpose**: Bump semantic versions (MAJOR.MINOR.PATCH) with intelligent defaults

**Usage**:
```bash
esh-cli bump-version <environment> [flags]
```

**Flags**:
- `--major`: Bump major version (breaking changes) - `1.2.3 → 2.0.0-1`
- `--minor`: Bump minor version (new features) - `1.2.3 → 1.3.0-1`
- `--patch`: Bump patch version (bug fixes) - `1.2.3 → 1.2.4-1`
- `--auto`: Auto-detect bump type from commit messages
- `--preview`: Show what would be created without executing
- `--service`: Target specific service
- `--from-commit`: Specify commit to tag (default: HEAD)

**Examples**:
```bash
# Manual version bumps
esh-cli bump-version stg6 --major     # Breaking changes
esh-cli bump-version stg6 --minor     # New features
esh-cli bump-version stg6 --patch     # Bug fixes

# Auto-detection based on conventional commits
esh-cli bump-version stg6 --auto      # Analyzes commit messages

# Preview mode (safe dry-run)
esh-cli bump-version stg6 --major --preview

# Service-specific tagging
esh-cli bump-version stg6 --patch --service api

# Tag specific commit
esh-cli bump-version stg6 --minor --from-commit abc123
```

**Conventional Commit Detection**:
- `feat:` → Minor bump
- `fix:` → Patch bump  
- `BREAKING CHANGE:` → Major bump
- `chore:`, `docs:`, `style:` → Patch bump

---

### `version-list` - Enhanced Version Listing

**Purpose**: List semantic versions with advanced filtering and formatting

**Usage**:
```bash
esh-cli version-list [environment] [flags]
```

**Flags**:
- `--all`: Show all environments
- `--major <n>`: Filter by major version
- `--minor <n>`: Filter by minor version  
- `--format <table|json|compact>`: Output format
- `--sort <version|date>`: Sort order
- `--limit <n>`: Maximum results per environment

**Examples**:
```bash
# Basic listing
esh-cli version-list stg6                    # All versions for stg6
esh-cli version-list --all                   # All environments

# Filtering
esh-cli version-list stg6 --major 1          # Major version 1.x.x
esh-cli version-list stg6 --major 1 --minor 2 # Version 1.2.x
esh-cli version-list stg6 --limit 5          # Last 5 versions

# Output formats
esh-cli version-list stg6 --format json      # JSON output
esh-cli version-list stg6 --format compact   # Compact view
esh-cli version-list stg6 --sort date        # Sort by date
```

**Output Formats**:
- **Table**: Human-readable with columns for tag, version, date, commit
- **JSON**: Machine-readable for scripting and automation
- **Compact**: Minimal output showing tag and version only

---

### `version-diff` - Version Comparison & Analysis

**Purpose**: Compare versions and analyze semantic differences

**Usage**:
```bash
esh-cli version-diff <tag1> [tag2] [flags]
esh-cli version-diff <environment> [flags]
```

**Flags**:
- `--history`: Show version history for environment
- `--remote`: Compare with remote tags
- `--commits`: Show commits between versions
- `--files`: Show changed files
- `--stats`: Show detailed statistics
- `--since <date>`: Show changes since date (YYYY-MM-DD)

**Examples**:
```bash
# Compare specific tags
esh-cli version-diff stg6_1.2.3-1 stg6_1.2.4-1

# Show environment history
esh-cli version-diff stg6 --history

# Show commits since tag
esh-cli version-diff stg6_1.2.3-1 --commits

# Detailed comparison with stats
esh-cli version-diff stg6_1.2.3-1 stg6_1.3.0-1 --stats --files

# Compare with previous tag (auto-detect)
esh-cli version-diff stg6_1.2.4-1    # Compares with previous version

# Show changes since date
esh-cli version-diff --since 2024-01-01
```

**Analysis Features**:
- Semantic version type detection (MAJOR/MINOR/PATCH)
- Commit count and contributors
- File change statistics
- Release timeline analysis

---

### `changelog` - Automated Changelog Generation

**Purpose**: Generate changelogs from git history and semantic versions

**Usage**:
```bash
esh-cli changelog [environment] [flags]
esh-cli changelog <tag1>..<tag2> [flags]
```

**Flags**:
- `--format <markdown|json|text>`: Output format
- `--conventional-commits`: Parse conventional commit messages
- `--full`: Generate complete changelog
- `--since <date>`: Changes since date (YYYY-MM-DD)
- `--output <file>`: Write to file instead of stdout
- `--from <tag>`: Start tag for range
- `--to <tag>`: End tag for range
- `--group-by-type`: Group entries by change type

**Examples**:
```bash
# Generate changelog for environment
esh-cli changelog stg6 --format markdown

# Between specific tags
esh-cli changelog stg6_1.2.0-1..stg6_1.3.0-1

# Parse conventional commits
esh-cli changelog --conventional-commits --group-by-type

# Full changelog to file
esh-cli changelog --full --output CHANGELOG.md

# Recent changes
esh-cli changelog --since 2024-01-01 --format json
```

**Conventional Commit Parsing**:
- Groups commits by type: feat, fix, docs, style, refactor, test, chore
- Extracts scope and breaking changes
- Formats according to conventional changelog standards

---

### `branch-version` - Git Flow Integration

**Purpose**: Branch-aware versioning and git flow integration

**Usage**:
```bash
esh-cli branch-version [flags]
```

**Flags**:
- `--suggest`: Suggest version bump based on current branch
- `--auto-tag`: Automatically create appropriate tag
- `--release-prep`: Prepare release branch workflow
- `--env <environment>`: Target environment for tagging
- `--service <service>`: Service name for tagging

**Examples**:
```bash
# Get branch-based suggestions
esh-cli branch-version --suggest

# Auto-create tag based on branch
esh-cli branch-version --auto-tag --env stg6

# Release preparation workflow
esh-cli branch-version --release-prep

# Service-specific branch tagging
esh-cli branch-version --auto-tag --env stg6 --service api
```

**Branch Naming Conventions**:
- `feature/*` → Minor version bump
- `hotfix/*`, `bugfix/*` → Patch version bump
- `release/*` → Release preparation workflow
- `develop`, `main` → Analyze commits for bump type
- `chore/*` → Patch version bump

---

## 🏷️ Traditional Tag Management Commands

### `add-tag` - Core Tag Management

**Purpose**: Add and manage traditional tags

**Usage**:
```bash
esh-cli add-tag <environment> <version|last> [flags]
```

**Flags**:
- `--from`: Tag to promote from
- `--hot-fix`: Tag hot fix (requires release branch)
- `--service`: Service name to tag

**Examples**:
```bash
# Show last tag
esh-cli add-tag stg6 last

# Add new tag
esh-cli add-tag stg6 1.2-1

# Promote between environments
esh-cli add-tag production2 1.2-1 --from stg6_1.2-0

# Hot fix tagging
esh-cli add-tag stg6 1.2-1 --hot-fix

# Service-specific tagging
esh-cli add-tag stg6 1.2-1 --service myservice
```

### `last-tag` - Query Last Tags

**Purpose**: Query the last tag for an environment

**Usage**:
```bash
esh-cli last-tag <environment> [flags]
```

**Examples**:
```bash
# Get last tag for environment
esh-cli last-tag stg6

# With service name
esh-cli last-tag stg6 --service api
```

---

## 🔧 Configuration & Management

### `projects` - Multi-Project Management

**Purpose**: Manage multiple project configurations

**Usage**:
```bash
esh-cli projects [flags]
```

### Global Flags

Available for all commands:
- `--config <file>`: Specify config file (default: $HOME/.esh-cli.yaml)
- `--help`: Show help for command
- `--version`: Show version information

---

## 📋 Environment & Tag Format

### Supported Environments
- `dev` - Development environment
- `mimic2` - Mimicry/testing environment
- `stg6` - Staging environment
- `demo` - Demo environment
- `production2` - Production environment

### Tag Format
Tags follow semantic versioning with environment prefixes:

**Format**: `[service_]env_major.minor.patch-release[.hotfix]`

**Examples**:
- `stg6_1.2.3-1` - Standard semantic version tag
- `stg6_1.2.3-1.1` - Hot fix tag
- `myservice_stg6_1.2.3-1` - Service-specific tag

### Semantic Version Rules
- **MAJOR**: Incompatible API changes (breaking changes)
- **MINOR**: Backward-compatible functionality additions  
- **PATCH**: Backward-compatible bug fixes
- **Release**: Environment-specific release number (starts at 1)
- **Hotfix**: Optional hotfix number for emergency fixes

---

## 🚀 Best Practices

### Version Bumping Strategy
1. **Use conventional commits** for automatic bump detection
2. **Preview changes** before applying with `--preview`
3. **Test in lower environments** before production
4. **Generate changelogs** for release documentation
5. **Use branch-based workflows** for consistent versioning

### Git Flow Integration
1. **Feature branches** → Minor version bumps
2. **Hotfix branches** → Patch version bumps
3. **Release branches** → Preparation workflows
4. **Main/develop** → Commit analysis based bumps

### Automation Tips
1. **CI/CD Integration**: Use in GitHub Actions/GitLab CI
2. **Preview Mode**: Always test with `--preview` first
3. **Service Isolation**: Use `--service` for microservices
4. **Cross-Environment**: Use traditional promotion for deployment chains
