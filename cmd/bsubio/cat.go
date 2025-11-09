package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
)

func runCat(args []string) error {
	fs := flag.NewFlagSet("cat", flag.ContinueOnError)

	// Custom usage function
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), "Usage: bsubio cat <jobid>\n\n")
		_, _ = fmt.Fprintf(fs.Output(), "Print job output (stdout)\n\n")
		_, _ = fmt.Fprintf(fs.Output(), "Arguments:\n")
		_, _ = fmt.Fprintf(fs.Output(), "  jobid    Job ID\n")
	}

	// Parse flags (none defined, but this handles help/errors)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Get remaining arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) != 1 {
		fs.Usage()
		return fmt.Errorf("expected 1 argument, got %d", len(remainingArgs))
	}

	jobID := remainingArgs[0]

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get job output
	resp, err := client.GetJobOutput(ctx, bsubio.JobId(jobID))
	if err != nil {
		return fmt.Errorf("failed to get job output: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get job output: HTTP %d", resp.StatusCode)
	}

	// Write output to stdout
	if _, err := os.Stdout.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
