# Scripts

## mkrel.sh - Release Script

Creates a new release of bsubio CLI with automated changelog generation and GitHub release.

### Prerequisites

- `git` - Version control
- `gh` - GitHub CLI ([install](https://cli.github.com/))
- `goreleaser` - Release automation ([install](https://goreleaser.com/install/))

Install missing tools:
```bash
brew install gh goreleaser
```

Authenticate with GitHub:
```bash
gh auth login
```

### Usage

From the repository root:

```bash
./scripts/mkrel.sh
```

The script will:
1. Check prerequisites (git, gh, goreleaser)
2. Verify working directory is clean
3. Display current branch and latest tag
4. Prompt for new version (e.g., v1.2.3)
5. Generate changelog from commits since last tag
6. Run `make check` to verify code quality
7. Create and push git tag
8. Run goreleaser to create GitHub release with binaries

### Version Format

Versions must follow semantic versioning with a `v` prefix:
- Release: `v1.2.3`
- Pre-release: `v1.2.3-beta.1`, `v1.2.3-rc.1`

### Changelog Generation

The script automatically generates changelogs by categorizing commits:
- **Features**: Commits starting with `feat:`, `Feat:`, or `FEAT:`
- **Bug Fixes**: Commits starting with `fix:`, `Fix:`, or `FIX:`
- **Other Changes**: All other commits

### Example

```bash
$ ./scripts/mkrel.sh

========================================
  bsubio CLI Release Script
========================================

➜ Checking prerequisites...
✓ All prerequisites met
➜ Checking working directory...
✓ Working directory is clean
➜ Current branch: main
➜ Latest tag: v0.1.0

Enter new version (e.g., v1.0.0): v0.2.0

## Changelog

### Features
- feat: Add device flow authentication
- feat: Add flag library support

### Bug Fixes
- fix: Handle HTTP client timeouts

⚠ Create release v0.2.0? (y/N) y
➜ Running make check...
✓ make check passed
➜ Creating tag v0.2.0...
➜ Pushing tag to origin...
✓ Tag pushed
➜ Running goreleaser...
✓ Release v0.2.0 created successfully!

➜ Release URL: https://github.com/bsubio/cli/releases/tag/v0.2.0

✓ 🎉 Release complete!
```

### Notes

- Always run from the `main` branch for production releases
- The script will warn if not on `main` but allows you to continue
- Failed `make check` will abort the release
- GoReleaser configuration is in `.goreleaser.yaml`
