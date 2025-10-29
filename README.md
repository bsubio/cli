# bsubio CLI

Command-line interface for [bsub.io](https://bsub.io) batch processing service.

## Installation


## Quick Start

First configure the CLI.

    bsubio config

The CLI will ask you for the API key.
You must have account at https://app.bsub.io first, to get it.
After you register, from "API Keys", click "New Key" and copy & paste the API key to
the terminal.
Configuration is stored in `~/.config/bsubio/config.json` and you can modify it at
any point in time.

Then you can do dry run:

    echo 123 > input.txt
    bsubio submit input.txt passthrough

This will submit the file `input.txt` that we just created with its  sample "123" string to `bsub.io` infrastructure.
The `passthrough` is like `cat` in command line: it should read input and print output without modification.

    bsubio status <jobid>
    bsubio cat <jobid>

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

## Support

For issues and questions:
- GitHub Issues: https://github.com/bsubio/cli/issues
- Documentation: https://docs.bsub.io
