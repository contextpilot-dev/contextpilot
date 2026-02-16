#!/bin/bash
# ContextPilot Release Script (Closed Source)
# Builds binaries and publishes to contextpilot.dev
#
# Usage: ./scripts/release-closed.sh v0.2.0

set -e

VERSION=$1
WEBSITE_REPO="${WEBSITE_REPO:-$HOME/Documents/source-code/contextpilot-website}"
HOMEBREW_TAP="${HOMEBREW_TAP:-/tmp/homebrew-tap}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { printf "${BLUE}==>${NC} %s\n" "$1"; }
success() { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn() { printf "${YELLOW}!${NC} %s\n" "$1"; }
error() { printf "${RED}✗${NC} %s\n" "$1" >&2; exit 1; }

# Validate version
if [ -z "$VERSION" ]; then
    echo "Usage: ./scripts/release-closed.sh <version>"
    echo "Example: ./scripts/release-closed.sh v0.2.0"
    exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-.*)?$ ]]; then
    error "Version must be in format v0.0.0 or v0.0.0-beta"
fi

NPM_VERSION=${VERSION#v}  # Remove 'v' prefix

echo ""
echo "  🦞 ContextPilot Release Script"
echo "  ================================"
echo "  Version: $VERSION"
echo ""

# Check repos exist
if [ ! -d "$WEBSITE_REPO" ]; then
    error "Website repo not found at $WEBSITE_REPO"
fi

# Clone homebrew tap if needed
if [ ! -d "$HOMEBREW_TAP" ]; then
    info "Cloning homebrew-tap..."
    git clone https://github.com/contextpilot-dev/homebrew-tap.git "$HOMEBREW_TAP"
fi

# Step 1: Build all platforms
info "Building binaries for all platforms..."
make dist VERSION=$VERSION
success "Built all binaries"

# Step 2: Create release directory in website repo
RELEASE_DIR="$WEBSITE_REPO/public/releases/$VERSION"
info "Creating release directory: $RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# Step 3: Copy binaries
info "Copying binaries..."
cp dist/*.tar.gz dist/*.zip dist/checksums.txt "$RELEASE_DIR/"
success "Copied binaries to website repo"

# Step 4: Update latest symlink
info "Updating 'latest' symlink..."
cd "$WEBSITE_REPO/public/releases"
rm -f latest
ln -sf "$VERSION" latest
success "Updated latest symlink"

# Step 5: Get checksums
info "Reading checksums..."
cd "$RELEASE_DIR"
SHA_DARWIN_AMD64=$(grep "darwin-amd64" checksums.txt | cut -d' ' -f1)
SHA_DARWIN_ARM64=$(grep "darwin-arm64" checksums.txt | cut -d' ' -f1)
SHA_LINUX_AMD64=$(grep "linux-amd64" checksums.txt | cut -d' ' -f1)
SHA_LINUX_ARM64=$(grep "linux-arm64" checksums.txt | cut -d' ' -f1)

echo "  darwin-amd64: $SHA_DARWIN_AMD64"
echo "  darwin-arm64: $SHA_DARWIN_ARM64"
echo "  linux-amd64:  $SHA_LINUX_AMD64"
echo "  linux-arm64:  $SHA_LINUX_ARM64"

# Step 6: Update Homebrew formula
info "Updating Homebrew formula..."
cat > "$WEBSITE_REPO/public/homebrew/contextpilot.rb" << EOF
# Homebrew formula for ContextPilot
# To use: brew tap contextpilot-dev/tap && brew install contextpilot

class Contextpilot < Formula
  desc "Make every AI tool understand your codebase"
  homepage "https://contextpilot.dev"
  version "$NPM_VERSION"
  license "MIT"

  on_macos do
    on_intel do
      url "https://contextpilot.dev/releases/v#{version}/contextpilot-darwin-amd64.tar.gz"
      sha256 "$SHA_DARWIN_AMD64"
    end
    on_arm do
      url "https://contextpilot.dev/releases/v#{version}/contextpilot-darwin-arm64.tar.gz"
      sha256 "$SHA_DARWIN_ARM64"
    end
  end

  on_linux do
    on_intel do
      url "https://contextpilot.dev/releases/v#{version}/contextpilot-linux-amd64.tar.gz"
      sha256 "$SHA_LINUX_AMD64"
    end
    on_arm do
      url "https://contextpilot.dev/releases/v#{version}/contextpilot-linux-arm64.tar.gz"
      sha256 "$SHA_LINUX_ARM64"
    end
  end

  def install
    binary_name = "contextpilot"
    if OS.mac?
      binary_name += "-darwin"
    else
      binary_name += "-linux"
    end
    binary_name += Hardware::CPU.arm? ? "-arm64" : "-amd64"
    
    bin.install binary_name => "contextpilot"
  end

  test do
    assert_match "contextpilot", shell_output("#{bin}/contextpilot --help")
  end
end
EOF
success "Updated Homebrew formula"

# Step 7: Copy formula to homebrew tap
info "Copying formula to homebrew-tap..."
cp "$WEBSITE_REPO/public/homebrew/contextpilot.rb" "$HOMEBREW_TAP/Formula/"
success "Copied formula"

# Step 8: Update npm package version
info "Updating npm package version..."
cd "$(dirname "$0")/.."
sed -i '' "s/\"version\": \".*\"/\"version\": \"$NPM_VERSION\"/" npm/package.json
success "Updated npm/package.json"

# Step 9: Commit and push website repo
info "Committing website repo..."
cd "$WEBSITE_REPO"
git add -A
git commit -m "Release $VERSION

- Add binaries for $VERSION
- Update Homebrew formula
- Update latest symlink"
git push origin main
success "Pushed website repo"

# Step 10: Commit and push homebrew tap
info "Committing homebrew-tap..."
cd "$HOMEBREW_TAP"
git add -A
git commit -m "Update to $VERSION"
git push origin main
success "Pushed homebrew-tap"

# Step 11: Tag the main repo (optional, even if private)
info "Creating git tag..."
cd "$(dirname "$0")/.."
git add npm/package.json
git commit -m "chore: bump version to $VERSION" || true
git tag -a "$VERSION" -m "Release $VERSION"
success "Created tag $VERSION"

echo ""
echo "=========================================="
success "Release $VERSION complete!"
echo "=========================================="
echo ""
echo "📦 Binaries published to:"
echo "   https://contextpilot.dev/releases/$VERSION/"
echo ""
echo "🍺 Homebrew updated:"
echo "   brew upgrade contextpilot"
echo ""
echo "📝 Next steps:"
echo "   1. Push tag: git push origin $VERSION"
echo "   2. Publish npm: cd npm && npm publish"
echo "   3. Test install: curl -fsSL https://contextpilot.dev/install.sh | sh"
echo ""
