# Organization Homebrew Tap Template

This directory contains templates for setting up an organization-wide Homebrew tap.

## Quick Setup

1. **Create a new repository** named `homebrew-tools` (or `homebrew-<org>-tools`) in your organization

2. **Initialize the repository structure**:
   ```bash
   mkdir Formula
   cp ../esh-cli-git/homebrew-formula/esh-cli.rb Formula/
   ```

3. **Create a README.md**:
   ```markdown
   # Organization Homebrew Tap
   
   Custom Homebrew formulas for [Your Organization].
   
   ## Installation
   
   ```bash
   # Add the tap
   brew tap PocketfulDev/tools
   
   # Install tools
   brew install esh-cli
   ```
   
   ## Available Formulas
   
   - `esh-cli` - Git tag management tool
   ```

4. **Commit and push**:
   ```bash
   git add .
   git commit -m "Initial tap setup with esh-cli formula"
   git push origin main
   ```

## Usage for Team Members

Once the tap is set up, team members can install tools with:

```bash
# One-time tap setup
brew tap PocketfulDev/tools

# Install any tool
brew install esh-cli

# Update tools
brew upgrade esh-cli
```

## Adding New Formulas

1. Create new formula files in the `Formula/` directory
2. Follow Homebrew formula conventions
3. Test locally with `brew install --build-from-source ./Formula/new-tool.rb`
4. Commit and push to make available organization-wide

## Automation

The `update-formula.sh` script in the main project automatically:
- Downloads release checksums
- Updates SHA256 hashes
- Updates version numbers

After running it, simply copy the updated formula to this tap repository.
