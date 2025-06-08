# ESH CLI Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Installation

#### Option 1: Homebrew (Recommended)
```bash
# Add the tap (one time)
brew tap PocketfulDev/tools

# Install ESH CLI
brew install esh-cli

# Verify installation
esh-cli --version
```

#### Option 2: Direct Download
```bash
# Download latest release for your platform
curl -L https://github.com/PocketfulDev/esh-cli/releases/latest/download/esh-cli-darwin-amd64 -o esh-cli
chmod +x esh-cli
./esh-cli --version
```

### Basic Usage

#### 1. Traditional Tag Management
```bash
# Show last tag for environment
esh-cli add-tag stg6 last

# Add new tag
esh-cli add-tag stg6 1.2-1

# Promote from staging to production
esh-cli add-tag production2 1.2-1 --from stg6_1.2-0
```

#### 2. Modern Semantic Versioning
```bash
# Intelligent version bumping
esh-cli bump-version stg6 --major     # Breaking changes
esh-cli bump-version stg6 --minor     # New features  
esh-cli bump-version stg6 --patch     # Bug fixes

# Auto-detect from commit messages
esh-cli bump-version stg6 --auto

# Preview changes safely
esh-cli bump-version stg6 --major --preview
```

#### 3. Version Analysis
```bash
# List versions with filtering
esh-cli version-list stg6
esh-cli version-list stg6 --major 1 --format json

# Compare versions
esh-cli version-diff stg6_1.2.3-1 stg6_1.2.4-1
esh-cli version-diff stg6 --history

# Generate changelogs
esh-cli changelog stg6 --conventional-commits
esh-cli changelog --full --output CHANGELOG.md
```

### Common Workflows

#### Development to Production Flow
```bash
# 1. Development
git checkout develop
echo "feat: new feature" > commit.txt
git commit -F commit.txt
esh-cli bump-version dev --auto              # Auto minor bump

# 2. Staging
git checkout release/1.3.0
esh-cli add-tag stg6 1.3.0-1 --from dev_1.3.0-1

# 3. Production (after testing)
git checkout main
esh-cli add-tag production2 1.3.0-1 --from stg6_1.3.0-1
esh-cli changelog production2 --output RELEASE_NOTES.md
```

#### Git Flow Integration
```bash
# Branch-based suggestions
git checkout feature/user-auth
esh-cli branch-version --suggest              # Suggests: MINOR

# Auto-create tags based on branch
esh-cli branch-version --auto-tag --env dev   # Creates appropriate tag

# Release preparation
git checkout release/1.4.0
esh-cli branch-version --release-prep         # Prepares release workflow
```

#### Service-Specific Tagging
```bash
# Tag specific microservice
esh-cli bump-version stg6 --patch --service api
esh-cli bump-version stg6 --minor --service frontend

# List service tags
esh-cli version-list stg6 --format json | jq '.[] | select(.service=="api")'
```

### Configuration

Create `~/.esh-cli.yaml` for defaults:
```yaml
# Default environment
default_environment: stg6

# Service configuration
services:
  api:
    path: "./services/api"
  frontend:
    path: "./services/frontend"

# Git configuration
git:
  auto_push: true
  confirm_before_push: true

# Conventional commits
conventional_commits:
  enabled: true
  types:
    feat: minor
    fix: patch
    breaking: major
```

### Best Practices

1. **Always preview first**: Use `--preview` flag for safety
2. **Use conventional commits**: Enable auto-detection
3. **Test in lower environments**: dev → staging → production
4. **Generate changelogs**: Document your releases
5. **Use branch-based workflows**: Consistent versioning

### Help & Support

```bash
# Command help
esh-cli --help
esh-cli bump-version --help

# Version information
esh-cli --version

# Check current environment
esh-cli add-tag stg6 last
```

### Next Steps

- **Read the [Command Reference](../reference/COMMAND_REFERENCE.md)** for detailed documentation
- **Explore [Testing Scenarios](../reference/TESTING_SCENARIOS.md)** for advanced usage
- **Check [Integrations](../reference/INTEGRATIONS.md)** for enterprise features
- **Review [GitHub Integration](../setup/GITHUB_TEST_INTEGRATION.md)** for CI/CD setup

---

🎉 **You're ready to start using ESH CLI!** Begin with traditional tag management, then graduate to semantic versioning as your team adopts the workflow.
