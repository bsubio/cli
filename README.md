# bsubio CLI

Command-line interface for [bsub.io](https://bsub.io) batch processing service.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/bsubio/cli.git
cd cli

make
```

The binary will be installed to `/usr/local/bin/bsubio`.

### Direct Build

```bash
go build -o bsubio ./cmd/bsubio
```

## Quick Start

1. Configure your API key:
```bash
bsubio config
```

2. Submit a job:
```bash
bsubio submit input.txt passthrough
```

3. Check job status:
```bash
bsubio status <jobid>
```

4. Get job output:
```bash
bsubio cat <jobid>
```

## Commands

### `bsubio config`

Configure API key and base URL. Configuration is stored securely in `~/.config/bsubio/config.json` with `0600` permissions.

```bash
bsubio config
```

### `bsubio submit`

Submit a job for processing.

```bash
bsubio submit [-o <file>] [-w] <input_file> <type>
```

**Flags:**
- `-w` - Wait for job to complete before returning
- `-o <file>` - Write output to file (requires `-w`)

**Examples:**
```bash
# Submit and get job ID
bsubio submit data.json json_format

# Submit, wait, and print output
bsubio submit -w input.txt passthrough

# Submit, wait, and save output to file
bsubio submit -w -o result.txt input.txt passthrough
```

### `bsubio wait`

Wait for a job to complete by polling its status.

```bash
bsubio wait [-v] [-t <seconds>] <jobid>
```

**Flags:**
- `-v` - Verbose output (show status updates)
- `-t <seconds>` - Polling interval in seconds (default: 5)

**Examples:**
```bash
# Wait with default interval
bsubio wait job_abc123

# Wait with verbose output and custom interval
bsubio wait -v -t 10 job_abc123
```

### `bsubio cat`

Print job output (stdout) to console.

```bash
bsubio cat <jobid>
```

**Example:**
```bash
bsubio cat job_abc123
```

### `bsubio logs`

Print job logs (stderr) to console.

```bash
bsubio logs <jobid>
```

**Example:**
```bash
bsubio logs job_abc123
```

### `bsubio jobs`

List recent jobs with optional filtering.

```bash
bsubio jobs [--status <status>] [--limit <n>]
```

**Flags:**
- `--status <status>` - Filter by job status (created, pending, claimed, finished, failed)
- `--limit <n>` - Maximum number of jobs to return (default: 20)

**Examples:**
```bash
# List recent jobs
bsubio jobs

# List only finished jobs
bsubio jobs --status finished

# List last 50 jobs
bsubio jobs --limit 50
```

### `bsubio status`

Show detailed information about a job.

```bash
bsubio status <jobid>
```

**Example:**
```bash
bsubio status job_abc123
```

Output includes:
- Job ID, type, and status
- Data size
- Timestamps (created, claimed, finished)
- Worker information
- Error messages (if failed)

### `bsubio cancel`

Cancel a job or all jobs.

```bash
bsubio cancel [-a|--all] <jobid>
```

**Flags:**
- `-a` or `--all` - Cancel all pending/claimed jobs

**Examples:**
```bash
# Cancel a single job
bsubio cancel job_abc123

# Cancel all jobs
bsubio cancel -a
```

### `bsubio rm`

Delete a job or all jobs.

```bash
bsubio rm [-a|--all] <jobid>
```

**Flags:**
- `-a` or `--all` - Delete all jobs

**Examples:**
```bash
# Delete a single job
bsubio rm job_abc123

# Delete all jobs
bsubio rm -a
```

### `bsubio types`

List available job processing types.

```bash
bsubio types
```

### `bsubio version`

Show CLI and API server version information.

```bash
bsubio version
```

### `bsubio help`

Show usage information and examples.

```bash
bsubio help
```

## Common Workflows

### Process a file and wait for results

```bash
bsubio submit -w -o result.txt input.txt passthrough
```

### Submit and monitor progress

```bash
# Submit job
JOBID=$(bsubio submit input.json json_format | grep -o 'Job submitted: .*' | cut -d' ' -f3)

# Wait for completion
bsubio wait -v $JOBID

# Get output
bsubio cat $JOBID
```

### Check recent jobs and their status

```bash
bsubio jobs --limit 10
```

### Debug a failed job

```bash
# Check status
bsubio status job_abc123

# View logs
bsubio logs job_abc123
```

## Configuration

Configuration is stored in `~/.config/bsubio/config.json`:

```json
{
  "api_key": "your-api-key-here",
  "base_url": "https://app.bsub.io"
}
```

## Development

### Project Structure

```
bsubio-cli/
├── cmd/bsubio/           # CLI implementation
│   ├── main.go          # Entry point and command router
│   ├── config.go        # Configuration management
│   ├── submit.go        # Submit command
│   ├── wait.go          # Wait command
│   ├── cat.go           # Cat command
│   ├── logs.go          # Logs command
│   ├── jobs.go          # Jobs command
│   ├── status.go        # Status command
│   ├── cancel.go        # Cancel command
│   ├── rm.go            # Delete command
│   ├── version.go       # Version command
│   └── types.go         # Types command
├── Makefile             # Build automation
└── go.mod               # Go dependencies
```

### Build Commands

```bash
# Build CLI binary
make build

# Clean build artifacts
make clean

# Install to system PATH
make install

# Download dependencies
make deps

# Run tests
make test

# Show help
make help
```

### Dependencies

- [bsubio-go](https://github.com/bsubio/bsubio-go) - Official Go SDK for bsub.io
- golang.org/x/term - Terminal password input

## Exit Codes

- `0` - Success
- `1` - Error (configuration, API, file system, etc.)

## License

See [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Support

For issues and questions:
- GitHub Issues: https://github.com/bsubio/cli/issues
- Documentation: https://docs.bsub.io
