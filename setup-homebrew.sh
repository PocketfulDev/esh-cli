#!/bin/bash

# setup-homebrew.sh - Complete setup guide for Homebrew distribution

set -e

echo "🍺 ESH CLI Homebrew Distribution Setup"
echo "====================================="
echo ""

# Check if we're in a git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

echo "✅ Git repository detected"

# Check if we have uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  Warning: You have uncommitted changes"
    echo "   These will be committed as part of the setup"
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled"
        exit 0
    fi
fi

echo ""
echo "📋 Setup Steps:"
echo "1. Push code to GitHub repository"
echo "2. Create and push first release tag"
echo "3. Set up organization Homebrew tap (optional)"
echo ""

# Get GitHub org/repo info
read -p "Enter your GitHub organization name: " GITHUB_ORG
read -p "Enter your repository name [esh-cli-git]: " REPO_NAME
REPO_NAME=${REPO_NAME:-esh-cli-git}

echo ""
echo "🔧 Updating configuration files..."

# Update Homebrew formula with correct GitHub org
sed -i.bak "s/your-org/$GITHUB_ORG/g" homebrew-formula/esh-cli.rb
sed -i.bak "s/your-org/$GITHUB_ORG/g" README.md
sed -i.bak "s/your-org/$GITHUB_ORG/g" HOMEBREW_SETUP.md

# Remove backup files
rm -f homebrew-formula/esh-cli.rb.bak README.md.bak HOMEBREW_SETUP.md.bak

# Commit any changes
if ! git diff-index --quiet HEAD --; then
    git add .
    git commit -m "Update GitHub organization references to $GITHUB_ORG"
fi

echo "✅ Configuration files updated"

# Check if remote exists
if ! git remote get-url origin > /dev/null 2>&1; then
    echo ""
    echo "🔗 Adding GitHub remote..."
    git remote add origin "https://github.com/$GITHUB_ORG/$REPO_NAME.git"
    echo "✅ Remote added: https://github.com/$GITHUB_ORG/$REPO_NAME.git"
else
    echo "✅ GitHub remote already configured"
fi

echo ""
echo "🚀 Next Steps:"
echo ""
echo "1. Create the GitHub repository:"
echo "   https://github.com/new"
echo "   Repository name: $REPO_NAME"
echo ""
echo "2. Push your code:"
echo "   git push -u origin main"
echo ""
echo "3. Create your first release:"
echo "   ./release.sh 1.0.0"
echo ""
echo "4. (Optional) Set up organization Homebrew tap:"
echo "   - Create repository: https://github.com/$GITHUB_ORG/homebrew-tools"
echo "   - See HOMEBREW_SETUP.md for detailed instructions"
echo ""
echo "5. Once release is complete, update formula:"
echo "   ./update-formula.sh 1.0.0 $GITHUB_ORG"
echo ""

echo "🎉 Setup complete! Follow the next steps above to publish your CLI tool."
