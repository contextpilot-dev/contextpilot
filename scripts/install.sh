#!/bin/sh
# ContextPilot installer
# Usage: curl -fsSL https://contextpilot.dev/install.sh | sh
#    or: curl -fsSL https://raw.githubusercontent.com/jitin-nhz/contextpilot/main/scripts/install.sh | sh

set -e

REPO="jitin-nhz/contextpilot"
BINARY="contextpilot"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1"
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*)  echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | \
        grep '"tag_name"' | \
        sed -E 's/.*"([^"]+)".*/\1/' || \
        error "Failed to fetch latest version. Check your internet connection."
}

# Download and install
install() {
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION="${VERSION:-$(get_latest_version)}"
    
    info "Installing ContextPilot ${VERSION} for ${OS}/${ARCH}..."
    
    # Construct download URL
    if [ "$OS" = "windows" ]; then
        FILENAME="${BINARY}-${OS}-${ARCH}.zip"
    else
        FILENAME="${BINARY}-${OS}-${ARCH}.tar.gz"
    fi
    
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT
    
    info "Downloading from ${DOWNLOAD_URL}..."
    
    if ! curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${FILENAME}"; then
        error "Failed to download ${FILENAME}"
    fi
    
    # Extract
    info "Extracting..."
    cd "${TMP_DIR}"
    
    if [ "$OS" = "windows" ]; then
        unzip -q "${FILENAME}"
        BINARY_FILE="${BINARY}-${OS}-${ARCH}.exe"
    else
        tar -xzf "${FILENAME}"
        BINARY_FILE="${BINARY}-${OS}-${ARCH}"
    fi
    
    # Install
    info "Installing to ${INSTALL_DIR}..."
    
    # Check if we need sudo
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${BINARY_FILE}" "${INSTALL_DIR}/${BINARY}"
        chmod +x "${INSTALL_DIR}/${BINARY}"
    else
        warn "Need elevated permissions to install to ${INSTALL_DIR}"
        sudo mv "${BINARY_FILE}" "${INSTALL_DIR}/${BINARY}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY}"
    fi
    
    # Verify installation
    if command -v "${BINARY}" >/dev/null 2>&1; then
        success "ContextPilot ${VERSION} installed successfully!"
        echo ""
        echo "Get started:"
        echo "  cd your-project"
        echo "  contextpilot init"
        echo ""
        echo "Documentation: https://contextpilot.dev"
    else
        warn "Installation complete, but ${BINARY} not found in PATH"
        echo "Add ${INSTALL_DIR} to your PATH, or run:"
        echo "  ${INSTALL_DIR}/${BINARY} --help"
    fi
}

# Check for required tools
check_deps() {
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
    fi
}

# Main
main() {
    echo ""
    echo "  ____            _            _   ____  _ _       _   "
    echo " / ___|___  _ __ | |_ _____  _| |_|  _ \(_) | ___ | |_ "
    echo "| |   / _ \| '_ \| __/ _ \ \/ / __| |_) | | |/ _ \| __|"
    echo "| |__| (_) | | | | ||  __/>  <| |_|  __/| | | (_) | |_ "
    echo " \____\___/|_| |_|\__\___/_/\_\\__|_|   |_|_|\___/ \__|"
    echo ""
    echo "Make every AI tool understand your codebase."
    echo ""
    
    check_deps
    install
}

main "$@"
