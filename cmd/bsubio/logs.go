package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func runLogs(args []string) error {
	fs := flag.NewFlagSet("logs", flag.ContinueOnError)

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio logs <jobid>\n\n")
		fmt.Fprintf(fs.Output(), "Show job logs (stderr)\n\n")
		fmt.Fprintf(fs.Output(), "Arguments:\n")
		fmt.Fprintf(fs.Output(), "  jobid    Job ID\n")
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

	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get job logs
	resp, err := client.GetJobLogs(ctx, jobUUID)
	if err != nil {
		return fmt.Errorf("failed to get job logs: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get job logs: HTTP %d", resp.StatusCode)
	}

	// Write logs to stdout
	if _, err := os.Stdout.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("failed to write logs: %w", err)
	}

	return nil
}
