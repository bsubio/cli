#!/usr/bin/env bash
set -e

# bsub.io CLI Installation Script
# Usage: curl -fsSL https://install.bsub.io/ | sh

# Configuration
REPO="bsubio/cli"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="bsubio"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
info() {
    echo -e "${GREEN}==>${NC} $1"
}

warn() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

error() {
    echo -e "${RED}Error:${NC} $1" >&2
    exit 1
}

# Detect OS
detect_os() {
    local os
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        *)          error "Unsupported operating system: $(uname -s)" ;;
    esac
    echo "$os"
}

# Detect architecture
detect_arch() {
    local arch
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        aarch64|arm64)  arch="arm64" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac
    echo "$arch"
}

# Get latest release version from GitHub
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | \
              grep '"tag_name":' | \
              sed -E 's/.*"tag_name": "v([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        error "Failed to get latest version from GitHub"
    fi

    echo "$version"
}

# Download binary
download_binary() {
    local os=$1
    local arch=$2
    local version=$3
    local binary_filename="${BINARY_NAME}-${os}-${arch}"
    local download_url="https://github.com/${REPO}/releases/download/v${version}/${binary_filename}"
    local checksum_url="${download_url}.sha256"

    info "Downloading bsubio v${version} for ${os}/${arch}..."

    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT

    # Download binary
    if ! curl -fsSL -o "${tmp_dir}/${binary_filename}" "${download_url}"; then
        error "Failed to download binary from ${download_url}"
    fi

    # Download checksum
    if ! curl -fsSL -o "${tmp_dir}/${binary_filename}.sha256" "${checksum_url}"; then
        warn "Failed to download checksum, skipping verification"
    else
        info "Verifying checksum..."
        (
            cd "${tmp_dir}"
            if command -v sha256sum >/dev/null 2>&1; then
                sha256sum -c "${binary_filename}.sha256" || error "Checksum verification failed"
            elif command -v shasum >/dev/null 2>&1; then
                shasum -a 256 -c "${binary_filename}.sha256" || error "Checksum verification failed"
            else
                warn "sha256sum or shasum not found, skipping checksum verification"
            fi
        )
    fi

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Install binary
    info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    mv "${tmp_dir}/${binary_filename}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    info "Installation complete!"
}

# Check if installation directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        warn "${INSTALL_DIR} is not in your PATH"
        echo ""
        echo "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo ""
        echo "    export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
    fi
}

# Main installation function
main() {
    info "bsub.io CLI Installer"
    echo ""

    # Detect system
    local os arch version
    os=$(detect_os)
    arch=$(detect_arch)
    info "Detected OS: ${os}"
    info "Detected Architecture: ${arch}"

    # Get latest version
    version=$(get_latest_version)
    info "Latest version: v${version}"
    echo ""

    # Download and install
    download_binary "$os" "$arch" "$version"

    # Check PATH
    echo ""
    check_path

    # Success message
    echo ""
    info "bsubio has been installed successfully!"
    echo ""
    echo "Run 'bsubio --help' to get started."
    echo "Run 'bsubio quickstart' for a quick tutorial."
}

# Run main function
main
