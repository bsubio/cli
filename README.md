# bsubio CLI

Command-line interface for [bsub.io](https://bsub.io) batch processing service.

## Installation

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
    $ bsubio submit input.txt passthrough

Expected output will be similar to:

    019a3256-26b4-7f1f-b1aa-0b45ab7b371d

Which is a JobID.

At any time you can see all jobs:

    $ bsubio jobs

    JOB ID                                   TYPE         STATUS          CREATED AT
    --------------------------------------------------------------------------------
    019a3256-26b4-7f1f-b1aa-0b45ab7b371d     passthrough  pending         2025-10-29 23:38

This will submit the file `input.txt` that we just created with its  sample "123" string to `bsub.io` infrastructure.
The `passthrough` is like `cat` in command line: it should read input and print output without modification.

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

## Support

For issues and questions:
- GitHub Issues: https://github.com/bsubio/cli/issues
- Documentation: https://docs.bsub.io
