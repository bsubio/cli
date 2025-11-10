# bsubio CLI

Command-line interface for [bsub.io](https://bsub.io) batch processing service.

## Installation

### Quick Install (Recommended)

For Linux and macOS:

```bash
curl -fsSL https://install.bsub.io/ | sh
```

The script will:
- Automatically detect your OS and architecture
- Download the latest release
- Verify checksums
- Install to `~/.local/bin/bsubio`
- Provide instructions to add to PATH if needed

### Manual Installation

Download the appropriate binary for your platform from the [latest release](https://github.com/bsubio/cli/releases/latest):

- **Linux (x86_64)**: `bsubio-linux-amd64`
- **Linux (ARM64)**: `bsubio-linux-arm64`
- **macOS (Intel)**: `bsubio-darwin-amd64`
- **macOS (Apple Silicon)**: `bsubio-darwin-arm64`
- **Windows (x86_64)**: `bsubio-windows-amd64.exe`

Then install manually:

```bash
# Download (replace with your platform)
curl -LO https://github.com/bsubio/cli/releases/latest/download/bsubio-linux-amd64

# Make executable
chmod +x bsubio-linux-amd64

# Move to your PATH
sudo mv bsubio-linux-amd64 /usr/local/bin/bsubio
```

### Build from Source

Requires Go 1.25 or later:

```bash
git clone https://github.com/bsubio/cli.git
cd cli
make build-static
sudo mv bin/bsubio /usr/local/bin/
```

## Quick Start

First configure the CLI.

    $ bsubio config

The CLI will ask you for the API key.
You must have account at https://app.bsub.io first, to get it.
After you register, from "API Keys", click "New Key" and copy & paste the API key to
the terminal.
Configuration is stored in `~/.config/bsubio/config.json` and you can modify it at
any point in time.

Then you can do dry run:

    $ echo 123 > input.txt
    $ bsubio submit input.txt passthru

Expected output will be similar to:

    019a3256-26b4-7f1f-b1aa-0b45ab7b371d

Which is a JobID.

At any time you can see all jobs:

    $ bsubio jobs

    JOB ID                                   TYPE         STATUS          CREATED AT
    --------------------------------------------------------------------------------
    019a3256-26b4-7f1f-b1aa-0b45ab7b371d     passthru     pending         2025-10-29 23:38

This will submit the file `input.txt` that we just created with its  sample "123" string to `bsub.io` infrastructure.
The `passthru` is like `cat` in command line: it should read input and print output without modification.

    bsubio status 019a3256-26b4-7f1f-b1aa-0b45ab7b371d
    bsubio cat 019a3256-26b4-7f1f-b1aa-0b45ab7b371d

## Exit Codes

- `0` - Success
- `1` - Error (configuration, API, file system, etc.)

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Development flow:

    git clone https://github.com/bsubio/cli.git
    cd cli
    make
    make clean
    make

For submitting new featurs:

1. Open PR to discuss the change
1. Fork the repository
1. Create a feature branch
1. Make your changes
1. Submit a pull request

### Release Process

This project uses **automated releases** based on [conventional commits](https://www.conventionalcommits.org/).

**How it works:**
1. Use conventional commit format in your commits:
   - `feat:` → triggers **minor** version bump (e.g., v0.1.0 → v0.2.0)
   - `fix:` → triggers **patch** version bump (e.g., v0.1.0 → v0.1.1)
2. When merged to `main`, GitHub Actions automatically:
   - Calculates the next version
   - Creates a git tag
   - Builds binaries for all platforms (Linux, macOS, Windows)
   - Generates changelog from commits
   - Publishes a GitHub release

**Commit examples:**
```bash
feat: add authentication support
fix: resolve memory leak in worker pool
```

**Note:** No manual tagging required! Releases happen automatically on merge to `main`.

## Support

For issues and questions:
- GitHub Issues: https://github.com/bsubio/cli/issues
- Documentation: https://docs.bsub.io

