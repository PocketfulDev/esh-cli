# ESH CLI

A Go CLI tool for managing git tags and deployments. This is a complete refactor of the original Python script using Cobra and Viper.

## Features

- Add and push hot fix tags
- Promote tags between environments
- Support for service-specific tags
- Git repository validation
- Interactive prompts for safety

## Installation

### Via Homebrew (Recommended)

```bash
# Add the tap (one-time setup)
brew tap eshos/tools

# Install esh-cli
brew install esh-cli

# Update to latest version
brew upgrade esh-cli
```

**Note**: Works with private repositories! See [PRIVATE_REPO_GUIDE.md](PRIVATE_REPO_GUIDE.md) for setup options.

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

### Basic Commands

Show help:
```bash
./esh-cli --help
./esh-cli add-tag --help
```

### Examples

Show last tag for staging:
```bash
./esh-cli add-tag stg6 ?
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

### Testing

```bash
make test
```

### Code formatting and vetting

```bash
make check
```

### Cleaning

```bash
make clean
```

## Releasing

### Creating a Release

1. **Prepare the release**:
   ```bash
   ./release.sh 1.0.0
   ```

2. **Update Homebrew formula** (after GitHub Actions completes):
   ```bash
   ./update-formula.sh 1.0.0 your-org
   ```

3. **Update your Homebrew tap** (if using organization tap):
   ```bash
   cp homebrew-formula/esh-cli.rb ../homebrew-your-org-tools/Formula/
   cd ../homebrew-your-org-tools
   git add Formula/esh-cli.rb
   git commit -m "Update esh-cli to v1.0.0"
   git push
   ```

For detailed setup instructions, see [HOMEBREW_SETUP.md](HOMEBREW_SETUP.md).

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
