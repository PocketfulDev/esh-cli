# Private Repository Homebrew Distribution Guide

## 🔒 Private Repository Options

### Option 1: Private Repo + Public Releases (Recommended)

**Best for**: Teams that want to keep source code private but don't mind public binaries

1. **Create private repository**
2. **Enable public releases** in repository settings:
   - Go to Settings → General → Features
   - Ensure "Releases" is enabled
   - Releases will be publicly accessible even with private repo

3. **Use the standard formula** (`esh-cli.rb`)
4. **Team installation**:
   ```bash
   brew tap PocketfulDev/tools
   brew install esh-cli
   ```

### Option 2: Fully Private Distribution

**Best for**: Maximum security - everything private

#### Setup for Fully Private:

1. **Use the private formula** (`esh-cli-private.rb`)
2. **Team members need GitHub tokens**:
   - Create at: https://github.com/settings/tokens
   - Required scope: `repo` (full repository access)

3. **Team installation process**:
   ```bash
   # Set GitHub token (add to ~/.zshrc for persistence)
   export HOMEBREW_GITHUB_API_TOKEN="ghp_your_token_here"
   
   # Install
   brew tap PocketfulDev/tools
   brew install esh-cli
   ```

#### Token Management for Teams:

**Option A: Individual Tokens**
- Each team member creates their own token
- Most secure but requires individual setup

**Option B: Shared Service Token**
- Create a service account with read-only access
- Share token with team (less secure but easier)

**Option C: GitHub App (Advanced)**
- Create GitHub App for your organization
- More complex setup but better for large teams

## 📋 Setup Instructions

### For Private Repo + Public Releases:
```bash
# Use standard setup
./setup-homebrew.sh
# When creating repo, choose "Private"
# Releases will still be public
```

### For Fully Private:
```bash
# 1. Run standard setup
./setup-homebrew.sh

# 2. Replace formula with private version
cp homebrew-formula/esh-cli-private.rb homebrew-formula/esh-cli.rb

# 3. Update your tap repository with the private formula
```

## 👥 Team Onboarding

### For Private Repo + Public Releases:
```bash
# One-time setup
brew tap PocketfulDev/tools
brew install esh-cli
```

### For Fully Private:
```bash
# 1. Create GitHub token at https://github.com/settings/tokens
#    Required scope: repo

# 2. Set environment variable (add to ~/.zshrc)
export HOMEBREW_GITHUB_API_TOKEN="your_token_here"

# 3. Install
brew tap PocketfulDev/tools
brew install esh-cli
```

## 🔧 Updating Private Formulas

The `update-formula.sh` script works with both approaches:

```bash
# For private repo + public releases
./update-formula.sh 1.0.1 your-org

# For fully private (same command)
./update-formula.sh 1.0.1 your-org
```

## 🚨 Security Considerations

1. **Tokens have access scope** - they can access other private repos
2. **Tokens can expire** - set appropriate expiration dates
3. **Audit token usage** - regularly review who has access
4. **Consider IP restrictions** if available in your GitHub plan

## 💡 Recommendations

- **Start with Option 1** (private repo + public releases)
- **Most teams find this sufficient** for internal tools
- **Only use fully private if** you have specific security requirements
- **Document the process** for your team in your internal wiki
