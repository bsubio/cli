#!/usr/bin/env bash
#
# Release script for bsubio CLI
# Uses gh, git, and goreleaser to create releases
#
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

info() {
    echo -e "${BLUE}➜ $1${NC}"
}

success() {
    echo -e "${GREEN}✓ $1${NC}"
}

warn() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    info "Checking prerequisites..."

    if ! command -v git &> /dev/null; then
        error "git is not installed"
    fi

    if ! command -v gh &> /dev/null; then
        error "gh (GitHub CLI) is not installed. Install with: brew install gh"
    fi

    if ! command -v goreleaser &> /dev/null; then
        error "goreleaser is not installed. Install with: brew install goreleaser"
    fi

    # Check if gh is authenticated
    if ! gh auth status &> /dev/null; then
        error "gh is not authenticated. Run: gh auth login"
    fi

    success "All prerequisites met"
}

# Check if working directory is clean
check_clean_working_dir() {
    info "Checking working directory..."

    if ! git diff-index --quiet HEAD --; then
        error "Working directory is not clean. Commit or stash changes first."
    fi

    success "Working directory is clean"
}

# Get current branch
get_current_branch() {
    git rev-parse --abbrev-ref HEAD
}

# Get latest tag
get_latest_tag() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
        error "Invalid version format: $version. Expected format: v1.2.3 or v1.2.3-beta.1"
    fi
}

# Generate changelog
generate_changelog() {
    local from_tag=$1
    local to_ref=${2:-HEAD}

    info "Generating changelog from ${from_tag} to ${to_ref}..."

    # Get commits since last tag
    local commits=$(git log ${from_tag}..${to_ref} --pretty=format:"%h %s" --no-merges)

    if [ -z "$commits" ]; then
        warn "No commits since last tag"
        return
    fi

    echo ""
    echo "## Changelog"
    echo ""

    # Parse commits by type
    local features=$(echo "$commits" | grep -E "^[a-f0-9]+ (feat|Feat|FEAT)" || true)
    local fixes=$(echo "$commits" | grep -E "^[a-f0-9]+ (fix|Fix|FIX)" || true)
    local other=$(echo "$commits" | grep -vE "^[a-f0-9]+ (feat|fix|Feat|Fix|FEAT|FIX)" || true)

    if [ -n "$features" ]; then
        echo "### Features"
        echo "$features" | sed 's/^[a-f0-9]* /- /'
        echo ""
    fi

    if [ -n "$fixes" ]; then
        echo "### Bug Fixes"
        echo "$fixes" | sed 's/^[a-f0-9]* /- /'
        echo ""
    fi

    if [ -n "$other" ]; then
        echo "### Other Changes"
        echo "$other" | sed 's/^[a-f0-9]* /- /'
        echo ""
    fi
}

# Main release process
main() {
    echo ""
    echo "========================================"
    echo "  bsubio CLI Release Script"
    echo "========================================"
    echo ""

    # Check prerequisites
    check_prerequisites

    # Check working directory
    check_clean_working_dir

    # Get current branch
    local current_branch=$(get_current_branch)
    info "Current branch: ${current_branch}"

    if [ "$current_branch" != "main" ]; then
        warn "Not on main branch. Continue anyway? (y/N)"
        read -r response
        if [[ ! $response =~ ^[Yy]$ ]]; then
            error "Aborted by user"
        fi
    fi

    # Get latest tag
    local latest_tag=$(get_latest_tag)
    info "Latest tag: ${latest_tag}"

    # Prompt for new version
    echo ""
    echo -n "Enter new version (e.g., v1.0.0): "
    read -r new_version

    # Validate version
    validate_version "$new_version"

    # Check if tag already exists
    if git rev-parse "$new_version" >/dev/null 2>&1; then
        error "Tag ${new_version} already exists"
    fi

    # Generate and display changelog
    echo ""
    local changelog=$(generate_changelog "$latest_tag")
    echo "$changelog"

    # Confirm release
    echo ""
    warn "Create release ${new_version}? (y/N)"
    read -r response
    if [[ ! $response =~ ^[Yy]$ ]]; then
        error "Aborted by user"
    fi

    # Run make check
    info "Running make check..."
    if ! make check; then
        error "make check failed. Fix issues before releasing."
    fi
    success "make check passed"

    # Create and push tag
    info "Creating tag ${new_version}..."
    git tag -a "$new_version" -m "Release ${new_version}"

    info "Pushing tag to origin..."
    git push origin "$new_version"
    success "Tag pushed"

    # Run goreleaser
    info "Running goreleaser..."
    if ! goreleaser release --clean; then
        error "goreleaser failed. Check the output above."
    fi

    success "Release ${new_version} created successfully!"

    # Get release URL
    local release_url=$(gh release view "$new_version" --json url -q .url 2>/dev/null || echo "")
    if [ -n "$release_url" ]; then
        echo ""
        info "Release URL: ${release_url}"
    fi

    echo ""
    success "🎉 Release complete!"
}

# Run main function
main "$@"
