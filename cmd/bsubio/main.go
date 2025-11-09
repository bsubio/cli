package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
)

var version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return runHelp(nil)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "register":
		return runRegister(args)
	case "config":
		return runConfig(args)
	case "submit":
		return runSubmit(args)
	case "wait":
		return runWait(args)
	case "cat":
		return runCat(args)
	case "logs":
		return runLogs(args)
	case "jobs":
		return runJobs(args)
	case "status":
		return runStatus(args)
	case "cancel":
		return runCancel(args)
	case "rm":
		return runRm(args)
	case "version":
		return runVersion(args)
	case "types":
		return runTypes(args)
	case "bench":
		return runBench(args)
	case "quickstart":
		return runQuickstart(args)
	case "help", "-h", "--help":
		return runHelpCommand(args)
	default:
		return fmt.Errorf("unknown command: %s\nRun 'bsubio help' for usage", command)
	}
}

func runHelp(args []string) error {
	fmt.Print(`bsubio - Command line tool for bsub.io batch processing

USAGE:
    bsubio <command> [options] [arguments]

COMMANDS:
    register                    Register with bsub.io using GitHub
    config                      Configure API key manually
    submit [-o <file>] [-w] <input_file> <type>
                                Submit a job for processing
    wait [-v] [-t <seconds>] <jobid>
                                Wait for a job to complete
    cat <jobid>                 Print job output (stdout)
    jobs [--status <status>] [--limit <n>]
                                List recent jobs
    status <jobid>              Show detailed job status
    logs <jobid>                Show job logs (stderr)
    cancel [-a|--all] <jobid>   Cancel a job (or all jobs with -a)
    rm [-a|--all] <jobid>       Delete a job (or all jobs with -a)
    version                     Show API server version
    types                       List available job types
    bench [options]             Benchmark job processing with test files
    quickstart                  Show quickstart guide
    help [command]              Show help message or help for a specific command

EXAMPLES:
    bsubio register
    bsubio config
    bsubio submit data.json json_format
    bsubio submit -w -o result.txt input.txt passthrough
    bsubio wait -v job_abc123
    bsubio cat job_abc123
    bsubio logs job_abc123
    bsubio status job_abc123
    bsubio cancel job_abc123
    bsubio cancel -a
    bsubio rm job_abc123
    bsubio rm -a
    bsubio jobs --limit 10
    bsubio types
    bsubio bench
    bsubio bench --type pdf_extract --dir tests/data
    bsubio version
    bsubio quickstart
    bsubio help submit
`)
	return nil
}

// createClient creates a new BSUB.IO client from config
func createClient() (*bsubio.BsubClient, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	client, err := bsubio.NewBsubClient(bsubio.Config{
		APIKey:  config.APIKey,
		BaseURL: config.BaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// getContext returns a context for API calls
func getContext() context.Context {
	return context.Background()
}
