# ESH CLI Project Structure

This document describes the organized structure of the ESH CLI project after cleanup and reorganization.

## 📁 Root Directory

```
esh-cli/
├── README.md                    # Main project documentation
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── main.go                      # Main application entry point
├── Makefile                     # Build automation
├── .gitignore                   # Git ignore patterns
├── .github/workflows/           # GitHub Actions CI/CD
├── build/                       # Build artifacts (gitignored)
├── dist/                        # Release distributions (gitignored)
├── docs/                        # Documentation hub
├── scripts/                     # Build and utility scripts
├── cmd/                         # CLI commands implementation
├── pkg/                         # Reusable packages
└── homebrew-*/                  # Homebrew distribution files
```

## 📚 Documentation Structure (`docs/`)

The documentation is organized into logical categories:

### `docs/guides/`
User-friendly tutorials and getting started materials:
- `QUICK_START_GUIDE.md` - 5-minute getting started guide

### `docs/reference/`
Complete technical documentation:
- `COMMAND_REFERENCE.md` - Full command documentation
- `TESTING_SCENARIOS.md` - Comprehensive testing strategies
- `INTEGRATIONS.md` - Enterprise integration examples

### `docs/setup/`
Installation and configuration guides:
- `GITHUB_TEST_INTEGRATION.md` - CI/CD setup
- `HOMEBREW_SETUP.md` - Package distribution
- `HOMEBREW_TAP_SETUP.md` - Custom tap configuration
- `ESHOS_TAP_SETUP.md` - ESH OS specific setup
- `PRIVATE_REPO_GUIDE.md` - Private repository configuration
- `DEPLOYMENT_CHECKLIST.md` - Production deployment

### `docs/design/`
Architecture and design documentation:
- `SEMANTIC_VERSIONING_DESIGN.md` - Design decisions and architecture

## 🔧 Build Structure (`build/` & `dist/`)

### `build/` (Development artifacts)
- `esh-cli*` - Development binaries
- `coverage.*` - Test coverage reports
- `test-results.json` - Test execution results

### `dist/` (Release artifacts)
- Platform-specific release binaries
- Generated during `make release-build`

## 🛠 Scripts (`scripts/`)

Utility and build scripts:
- `release.sh` - Release automation
- `update-formula.sh` - Homebrew formula updates
- `setup-homebrew.sh` - Homebrew setup automation
- `test-github-integration.sh` - Integration testing
- `get-coverage.sh` - Coverage calculation

## 💻 Source Code (`cmd/` & `pkg/`)

### `cmd/` - CLI Commands
Each command is implemented as a separate file:
- `root.go` - Root command and CLI setup
- `add-tag.go` - Traditional tag management
- `bump-version.go` - Semantic version bumping
- `version-list.go` - Version listing and filtering
- `version-diff.go` - Version comparison
- `changelog.go` - Changelog generation
- `branch-version.go` - Git flow integration
- `init.go` - Project initialization
- `last-tag.go` - Tag querying
- `projects.go` - Project management

### `pkg/utils/` - Core Utilities
- `utils.go` - General utilities and git operations
- `semver.go` - Semantic versioning logic
- `*_test.go` - Comprehensive test suites

## 🏗 Build System

### Makefile Targets
- `make build` - Build development binary to `build/`
- `make release-build` - Build all platform binaries to `dist/`
- `make test` - Run test suite
- `make test-coverage` - Generate coverage reports in `build/`
- `make clean` - Clean build artifacts

### GitHub Actions
- `ci.yml` - Continuous integration testing
- `release.yml` - Automated releases
- `badge-update.yml` - Dynamic README badges

## 📦 Distribution

### Homebrew Distribution
- `homebrew-formula/` - Formula definitions
- `homebrew-tap-template/` - Template for custom taps

## 🎯 Benefits of This Organization

1. **Clear separation of concerns** - Documentation, code, scripts, and build artifacts are properly separated
2. **Improved discoverability** - Logical documentation hierarchy with clear navigation
3. **Clean development environment** - Build artifacts contained in dedicated directories
4. **Maintainable CI/CD** - Consistent paths across all automation
5. **Professional structure** - Follows Go and open-source project conventions
6. **Easy onboarding** - New contributors can quickly understand the project layout

## 🔄 Migration Impact

All file references have been updated in:
- ✅ README.md links to documentation
- ✅ GitHub Actions workflows
- ✅ Makefile build targets
- ✅ Documentation cross-references
- ✅ Git ignore patterns

The reorganization maintains full backward compatibility for end users while providing a much cleaner development experience.
