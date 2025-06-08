# Homebrew Distribution Setup Guide

This guide explains how to distribute your `esh-cli` tool via Homebrew for easy installation across your organization.

## Prerequisites

1. **GitHub Repository**: Your code must be in a public GitHub repository
2. **GitHub Releases**: You'll create releases with pre-built binaries
3. **Homebrew Tap** (Optional): For organization-wide distribution

## Step 1: Prepare Your Repository

1. **Push your code to GitHub**:
   ```bash
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/your-org/esh-cli-git.git
   git push -u origin main
   ```

2. **Create and push your first tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Step 2: GitHub Actions Release

The GitHub Actions workflow (`.github/workflows/release.yml`) will automatically:
- Build binaries for macOS (Intel & Apple Silicon) and Linux
- Create compressed archives
- Generate SHA256 checksums
- Create a GitHub release

When you push a tag (e.g., `v1.0.0`), the release will be created automatically.

## Step 3: Create a Homebrew Tap (Organization-wide)

For organization-wide distribution, create a separate repository for your Homebrew tap:

1. **Create a new repository**: `homebrew-your-org-tools` (or similar)

2. **Create the tap structure**:
   ```
   homebrew-your-org-tools/
   └── Formula/
       └── esh-cli.rb
   ```

3. **Copy the formula**: Use the generated `homebrew-formula/esh-cli.rb` file

## Step 4: Update Formula After Each Release

After creating a new release:

1. **Update the formula with checksums**:
   ```bash
   ./update-formula.sh 1.0.0 your-org
   ```

2. **Copy updated formula to your tap repository**:
   ```bash
   cp homebrew-formula/esh-cli.rb ../homebrew-your-org-tools/Formula/
   ```

3. **Commit and push the tap**:
   ```bash
   cd ../homebrew-your-org-tools
   git add Formula/esh-cli.rb
   git commit -m "Update esh-cli to v1.0.0"
   git push
   ```

## Step 5: Installation Instructions

### For Organization Members

1. **Add the tap**:
   ```bash
   brew tap your-org/your-org-tools
   ```

2. **Install the CLI**:
   ```bash
   brew install esh-cli
   ```

3. **Update to latest version**:
   ```bash
   brew upgrade esh-cli
   ```

### Alternative: Direct Installation (without tap)

Users can also install directly from the formula URL:
```bash
brew install https://raw.githubusercontent.com/your-org/homebrew-your-org-tools/main/Formula/esh-cli.rb
```

## Automation Script

Here's a complete release script you can use:

```bash
#!/bin/bash
# release.sh - Automate the release process

set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.1"
    exit 1
fi

echo "Creating release $VERSION..."

# Create and push tag
git tag "v$VERSION"
git push origin "v$VERSION"

echo "Release v$VERSION created! GitHub Actions will build the binaries."
echo "Once the release is complete, run:"
echo "  ./update-formula.sh $VERSION your-org"
```

## Directory Structure

Your final project structure should look like:
```
esh-cli-git/
├── .github/workflows/release.yml
├── homebrew-formula/esh-cli.rb
├── update-formula.sh
├── release.sh
├── cmd/
├── pkg/
├── main.go
├── Makefile
└── README.md
```

## Security Considerations

- **Private Repositories**: If your code is private, consider using GitHub Packages or internal package repositories
- **Access Control**: Use Homebrew taps in private repositories for organization-only access
- **Binary Verification**: The SHA256 checksums ensure binary integrity

## Troubleshooting

1. **Release not triggered**: Check that your tag follows the pattern `v*` (e.g., `v1.0.0`)
2. **Build failures**: Check the GitHub Actions logs in your repository
3. **Formula errors**: Test the formula locally with `brew install --build-from-source`
4. **Checksum mismatches**: Re-run `update-formula.sh` with the correct version and organization
