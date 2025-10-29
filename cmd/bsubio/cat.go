package main

import (
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
)

func runCat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bsubio cat <jobid>")
	}

	jobID := args[0]

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get job output: HTTP %d", resp.StatusCode)
	}

	// Write output to stdout
	if _, err := os.Stdout.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
