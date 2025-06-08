# ESH CLI

[![Build Status](https://github.com/PocketfulDev/esh-cli/workflows/CI/badge.svg)](https://github.com/PocketfulDev/esh-cli/actions)
[![Test Coverage](https://img.shields.io/badge/coverage-32.3%25-orange)](https://github.com/PocketfulDev/esh-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/PocketfulDev/esh-cli)](https://goreportcard.com/report/github.com/PocketfulDev/esh-cli)

A Go CLI tool for managing git tags and deployments. This is a complete refactor of the original Python script using Cobra and Viper.

## 📚 Documentation

**[📖 Complete Documentation](docs/README.md)** - Browse all guides, references, and setup instructions

### Quick Links
- **[🚀 Quick Start Guide](docs/guides/QUICK_START_GUIDE.md)** - Get started in 5 minutes
- **[📖 Command Reference](docs/reference/COMMAND_REFERENCE.md)** - Complete command documentation
- **[⚙️ Setup Guides](docs/setup/)** - Installation and deployment guides
- **[🔗 Integrations](docs/reference/INTEGRATIONS.md)** - Enterprise integrations

## Features

### Core Functionality
- Add and push hot fix tags
- Promote tags between environments  
- Support for service-specific tags
- Git repository validation
- Interactive prompts for safety

### 🚀 Semantic Versioning (NEW)
- **Intelligent version bumping** with `--major`, `--minor`, `--patch` flags
- **Auto-detection** from conventional commit messages
- **Branch-aware versioning** with git flow integration
- **Version comparison** and semantic difference analysis
- **Automated changelog generation** from git history
- **Cross-environment version tracking** and synchronization
- **Preview mode** for safe version management

### Infrastructure
- **Comprehensive test coverage** (152+ tests) with automated reporting
- **Cross-platform builds** (macOS, Linux, ARM64 support)
- **CI/CD integration** ready with GitHub Actions
- **Enterprise-grade** semantic versioning capabilities

## Installation

### Via Homebrew (Recommended)

```bash
# Add the tap (one-time setup)
brew tap PocketfulDev/tools

# Install esh-cli
brew install esh-cli

# Update to latest version
brew upgrade esh-cli
```

**Note**: Works with private repositories! See [PRIVATE_REPO_GUIDE.md](docs/setup/PRIVATE_REPO_GUIDE.md) for setup options.

### Build from source

```bash
git clone https://github.com/PocketfulDev/esh-cli.git
cd esh-cli
make build
```

### Install to GOPATH

```bash
make install
```

### Download pre-built binaries

Download the latest release from the [releases page](https://github.com/PocketfulDev/esh-cli/releases).

## Usage

### 🏷️ Traditional Tag Management

Show help:
```bash
./esh-cli --help
./esh-cli add-tag --help
```

Show last tag for staging:
```bash
./esh-cli add-tag stg6 last
```

Add tag for staging on latest commit:
```bash
./esh-cli add-tag stg6 1.2-1
```

Promote from staging to production:
```bash
./esh-cli add-tag production2 1.2-1 --from stg6_1.2-0
```

Add hot fix tag (must be on release branch):
```bash
./esh-cli add-tag stg6 1.2-1 --hot-fix
```

Add tag with service name:
```bash
./esh-cli add-tag stg6 1.2-1 --service myservice
```

### 🚀 Semantic Versioning Commands

#### Version Bumping
```bash
# Bump semantic versions with flags
./esh-cli bump-version stg6 --major     # 1.2.3 → 2.0.0-1
./esh-cli bump-version stg6 --minor     # 1.2.3 → 1.3.0-1  
./esh-cli bump-version stg6 --patch     # 1.2.3 → 1.2.4-1

# Auto-detect from conventional commits
./esh-cli bump-version stg6 --auto      # Analyzes commit messages

# Preview mode (safe dry-run)
./esh-cli bump-version stg6 --major --preview

# Service-specific bumping
./esh-cli bump-version stg6 --patch --service myservice
```

#### Version Listing & Filtering
```bash
# List versions with advanced filtering
./esh-cli version-list stg6                    # All versions for stg6
./esh-cli version-list stg6 --major 1          # Filter by major version
./esh-cli version-list stg6 --format json      # JSON output
./esh-cli version-list --all                   # Cross-environment view
```

#### Version Comparison & Analysis
```bash
# Compare specific versions
./esh-cli version-diff stg6_1.2.3-1 stg6_1.2.4-1

# Show version history
./esh-cli version-diff stg6 --history

# Show commits between versions
./esh-cli version-diff stg6_1.2.3-1 --commits

# Detailed statistics
./esh-cli version-diff stg6_1.2.3-1 stg6_1.2.4-1 --stats
```

#### Changelog Generation
```bash
# Generate markdown changelog
./esh-cli changelog stg6 --format markdown

# Parse conventional commits
./esh-cli changelog --conventional-commits

# Full changelog to file
./esh-cli changelog --full --output CHANGELOG.md

# Between specific tags
./esh-cli changelog stg6_1.2.0-1..stg6_1.3.0-1
```

#### Git Flow Integration
```bash
# Branch-based version suggestions
./esh-cli branch-version --suggest

# Auto-create tag based on branch type
./esh-cli branch-version --auto-tag stg6

# Release preparation workflow
./esh-cli branch-version --release-prep
```

### Traditional Examples
```bash
./esh-cli add-tag stg6 1.2-1 --service myservice
```

## Tag Format

Tags follow the format: `[service_]env_major.minor.patch-release[.hotfix]`

Examples:
- `stg6_1.2-0` - Standard tag
- `stg6_1.2-0.1` - Hot fix tag
- `myservice_stg6_1.2-0` - Service-specific tag

## Supported Environments

- `dev`
- `mimic2`
- `stg6`
- `demo`
- `production2`

## Flags

- `-f, --from`: Tag to promote from
- `--hot-fix`: Tag hot fix (requires release branch)
- `-s, --service`: Service name to tag
- `--config`: Config file (default is $HOME/.esh-cli.yaml)

## Development

### Building

```bash
make build
```

## Development

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# View coverage in browser
make test-coverage && go tool cover -html=coverage.out
```

### Code Quality

```bash
# Run formatting and linting checks
make check

# Fix formatting issues
go fmt ./...
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## 🧪 Test Integration

This project features comprehensive test integration with GitHub:

- **Automated Coverage Reports**: Posted to pull requests
- **Test Result Visualization**: In GitHub Actions UI
- **Quality Gates**: Tests must pass before releases
- **Cross-Platform Testing**: Linux, macOS (Intel & ARM)

See [GITHUB_TEST_INTEGRATION.md](docs/setup/GITHUB_TEST_INTEGRATION.md) for detailed information about test reporting and coverage.

### Coverage Thresholds
- **Overall**: 30% (CLI tools have integration-heavy paths)
- **Utils Package**: 60% (core business logic)

## Releasing

### Creating a Release

1. **Prepare the release**:
   ```bash
   ./scripts/release.sh 1.0.0
   ```

2. **Update Homebrew formula** (after GitHub Actions completes):
   ```bash
   ./scripts/update-formula.sh 1.0.0 your-org
   ```

3. **Update your Homebrew tap** (if using organization tap):
   ```bash
   cp homebrew-formula/esh-cli.rb ../homebrew-your-org-tools/Formula/
   cd ../homebrew-your-org-tools
   git add Formula/esh-cli.rb
   git commit -m "Update esh-cli to v1.0.0"
   git push
   ```

For detailed setup instructions, see [HOMEBREW_SETUP.md](docs/setup/HOMEBREW_SETUP.md).

## Migration from Python

This Go version maintains full compatibility with the original Python script:

| Python Flag | Go Flag | Description |
|-------------|---------|-------------|
| `-f, --from` | `-f, --from` | Tag to promote from |
| `-hf, --hot_fix` | `--hot-fix` | Tag hot fix |
| `-s, --service` | `-s, --service` | Service name to tag |

The command structure remains the same:
- `python add_tag_hf.py stg6 1.2-1` → `./esh-cli add-tag stg6 1.2-1`
- All validation logic and behavior is preserved

## Configuration

You can create a configuration file at `$HOME/.esh-cli.yaml` for default settings (using Viper for configuration management).
