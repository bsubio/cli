package main

import (
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
)

func runLogs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bsubio logs <jobid>")
	}

	jobID := args[0]

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get job logs
	resp, err := client.GetJobLogs(ctx, bsubio.JobId(jobID))
	if err != nil {
		return fmt.Errorf("failed to get job logs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get job logs: HTTP %d", resp.StatusCode)
	}

	// Write logs to stdout
	if _, err := os.Stdout.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("failed to write logs: %w", err)
	}

	return nil
}
