# Setting Up PocketfulDev/homebrew-tools Tap

This guide will walk you through setting up the `PocketfulDev/homebrew-tools` Homebrew tap for distributing the ESH CLI tool.

## Overview

**Repository Structure:**
- **Source Code**: `PocketfulDev/esh-cli` (stays as is)
- **Homebrew Tap**: `PocketfulDev/homebrew-tools` (new repository to create)
- **User Installation**: `brew tap PocketfulDev/tools && brew install esh-cli`

## Step 1: Create the Homebrew Tap Repository

1. **Go to GitHub** and create a new repository:
   - **Organization**: `PocketfulDev`
   - **Repository name**: `homebrew-tools`
   - **Description**: "Homebrew formulae for ESH CLI tools"
   - **Visibility**: Public (Homebrew taps must be public)
   - **Initialize**: Yes, with README

2. **Repository URL**: `https://github.com/PocketfulDev/homebrew-tools`

## Step 2: Clone and Set Up the Tap Repository

```bash
# Clone the new tap repository
git clone https://github.com/PocketfulDev/homebrew-tools.git
cd homebrew-tools

# Create the Formula directory (required for Homebrew taps)
mkdir -p Formula

# Copy the esh-cli formula from your source project
cp /Users/jonathanpick/esh-cli-git/homebrew-formula/esh-cli.rb Formula/esh-cli.rb

# Create a basic README for the tap
cat << 'EOF' > README.md
# PocketfulDev Homebrew Tap

Custom Homebrew formulas for PocketfulDev tools.

## Installation

```bash
# Add the tap
brew tap PocketfulDev/tools

# Install tools
brew install esh-cli
```

## Available Formulas

- **esh-cli** - Git tag management and deployment tool

## Usage

Once installed, you can use the ESH CLI tool:

```bash
# Show help
esh-cli --help

# Add a tag
esh-cli add-tag staging 1.2.3

# Show version
esh-cli --version
```

## Updating

```bash
# Update the tap and upgrade tools
brew update
brew upgrade esh-cli
```
EOF

# Commit and push the initial setup
git add .
git commit -m "Initial tap setup with esh-cli formula"
git push origin main
```

## Step 3: Create Your First Release

Before the Homebrew formula will work, you need to create a release of the ESH CLI:

```bash
# Go back to your ESH CLI project
cd /Users/jonathanpick/esh-cli-git

# Create and push the first release tag
git tag v1.0.0
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow that:
1. Runs all tests
2. Builds cross-platform binaries
3. Creates a GitHub release with downloadable artifacts

## Step 4: Update the Formula with Real Checksums

After the release is created, you need to update the SHA256 checksums in the formula:

```bash
# Use the update script to get real checksums
./update-formula.sh 1.0.0

# Copy the updated formula to your tap
cp homebrew-formula/esh-cli.rb /path/to/homebrew-tools/Formula/esh-cli.rb

# Or manually update the SHA256 values by downloading the release artifacts:
# curl -L https://github.com/PocketfulDev/esh-cli/releases/download/v1.0.0/esh-cli-darwin-arm64.tar.gz | sha256sum
```

Then update the tap repository:

```bash
cd /path/to/homebrew-tools
git add Formula/esh-cli.rb
git commit -m "Update esh-cli formula with v1.0.0 checksums"
git push origin main
```

## Step 5: Test the Installation

```bash
# Test the complete workflow
brew tap PocketfulDev/tools
brew install esh-cli

# Verify it works
esh-cli --version
esh-cli --help
```

## Directory Structure

After setup, your tap repository should look like:

```
homebrew-tools/
├── README.md
└── Formula/
    └── esh-cli.rb
```

## Automated Updates

For future releases, use the provided automation:

```bash
# In your ESH CLI project directory
# 1. Create new release
git tag v1.0.1
git push origin v1.0.1

# 2. Wait for GitHub Actions to complete

# 3. Update formula automatically
./update-formula.sh 1.0.1

# 4. Copy to tap repository
cp homebrew-formula/esh-cli.rb /path/to/homebrew-tools/Formula/esh-cli.rb

# 5. Commit to tap
cd /path/to/homebrew-tools
git add Formula/esh-cli.rb
git commit -m "Update esh-cli to v1.0.1"
git push origin main
```

## Team Usage

Once set up, team members can install the tool with:

```bash
# One-time setup
brew tap PocketfulDev/tools

# Install
brew install esh-cli

# Update
brew upgrade esh-cli

# Uninstall (if needed)
brew uninstall esh-cli
brew untap PocketfulDev/tools
```

## Troubleshooting

### Formula Not Found
```bash
# Ensure tap is added correctly
brew tap PocketfulDev/tools

# Update Homebrew
brew update

# Try again
brew install esh-cli
```

### SHA256 Mismatch
```bash
# This means the checksums in the formula don't match the release artifacts
# Re-run the update script:
./update-formula.sh <version>
```

### Permission Issues
```bash
# Ensure you have write access to the PocketfulDev/homebrew-tools repository
# Check your GitHub permissions
```

## Next Steps

1. **Create the PocketfulDev/homebrew-tools repository** on GitHub
2. **Follow the setup steps above** to initialize it
3. **Create your first ESH CLI release** (v1.0.0)
4. **Update the formula** with real checksums
5. **Test the installation** process

After this setup, your team will be able to install the ESH CLI tool with a simple `brew install esh-cli` command after adding the tap once.
