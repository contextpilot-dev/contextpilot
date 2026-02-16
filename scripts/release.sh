#!/bin/bash
# Release script for ContextPilot
# Usage: ./scripts/release.sh v0.1.0

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: ./scripts/release.sh <version>"
    echo "Example: ./scripts/release.sh v0.1.0"
    exit 1
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-.*)?$ ]]; then
    echo "Error: Version must be in format v0.0.0 or v0.0.0-beta"
    exit 1
fi

echo "🚀 Releasing ContextPilot $VERSION"
echo ""

# Update npm package.json version
NPM_VERSION=${VERSION#v}  # Remove 'v' prefix
sed -i '' "s/\"version\": \".*\"/\"version\": \"$NPM_VERSION\"/" npm/package.json
echo "✓ Updated npm/package.json to $NPM_VERSION"

# Update Homebrew formula version
sed -i '' "s/version \".*\"/version \"$NPM_VERSION\"/" scripts/homebrew/contextpilot.rb
echo "✓ Updated Homebrew formula to $NPM_VERSION"

# Build all platforms
echo ""
echo "Building binaries..."
make dist VERSION=$VERSION

# Show checksums (need to update Homebrew formula manually after release)
echo ""
echo "📋 SHA256 Checksums (for Homebrew formula):"
cat dist/checksums.txt
echo ""

# Commit version changes
git add npm/package.json scripts/homebrew/contextpilot.rb
git commit -m "chore: bump version to $VERSION"
echo "✓ Committed version bump"

# Create and push tag
git tag -a "$VERSION" -m "Release $VERSION"
echo "✓ Created tag $VERSION"

echo ""
echo "Ready to release! Run:"
echo "  git push origin main"
echo "  git push origin $VERSION"
echo ""
echo "GitHub Actions will create the release automatically."
echo ""
echo "After release, update Homebrew formula SHA256 hashes from:"
echo "  https://github.com/jitin-nhz/contextpilot/releases/download/$VERSION/checksums.txt"
