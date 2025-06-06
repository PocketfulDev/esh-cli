#!/bin/bash

# Script to update Homebrew formula with release checksums
# Usage: ./update-formula.sh <version> <github-org>

set -e

VERSION=${1:-"1.0.0"}
GITHUB_ORG=${2:-"your-org"}
FORMULA_FILE="homebrew-formula/esh-cli.rb"

if [ ! -f "$FORMULA_FILE" ]; then
    echo "Error: Formula file $FORMULA_FILE not found"
    exit 1
fi

echo "Updating formula for version $VERSION..."

# Function to get SHA256 from GitHub release
get_sha256() {
    local filename=$1
    local url="https://github.com/$GITHUB_ORG/esh-cli-git/releases/download/v$VERSION/checksums.txt"
    
    echo "Fetching checksums from: $url"
    curl -sL "$url" | grep "$filename" | cut -d' ' -f1
}

# Get checksums for each platform
DARWIN_ARM64_SHA=$(get_sha256 "esh-cli-darwin-arm64.tar.gz")
DARWIN_AMD64_SHA=$(get_sha256 "esh-cli-darwin-amd64.tar.gz")
LINUX_ARM64_SHA=$(get_sha256 "esh-cli-linux-arm64.tar.gz")
LINUX_AMD64_SHA=$(get_sha256 "esh-cli-linux-amd64.tar.gz")

echo "Got checksums:"
echo "  Darwin ARM64: $DARWIN_ARM64_SHA"
echo "  Darwin AMD64: $DARWIN_AMD64_SHA"
echo "  Linux ARM64:  $LINUX_ARM64_SHA"
echo "  Linux AMD64:  $LINUX_AMD64_SHA"

# Update the formula file
sed -i.bak \
    -e "s/version \".*\"/version \"$VERSION\"/" \
    -e "s|https://github.com/your-org/|https://github.com/$GITHUB_ORG/|g" \
    -e "s/REPLACE_WITH_ARM64_SHA256/$DARWIN_ARM64_SHA/" \
    -e "s/REPLACE_WITH_AMD64_SHA256/$DARWIN_AMD64_SHA/" \
    -e "s/REPLACE_WITH_LINUX_ARM64_SHA256/$LINUX_ARM64_SHA/" \
    -e "s/REPLACE_WITH_LINUX_AMD64_SHA256/$LINUX_AMD64_SHA/" \
    "$FORMULA_FILE"

echo "Formula updated successfully!"
echo "Please review the changes in $FORMULA_FILE"
