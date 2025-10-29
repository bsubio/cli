package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bsubio/bsubio-go"
)

func runWait(args []string) error {
	var (
		verbose  bool
		interval int = 5 // default 5 seconds
		jobID    string
	)

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-v" {
			verbose = true
			i++
		} else if arg == "-t" {
			if i+1 >= len(args) {
				return fmt.Errorf("-t flag requires a timeout value in seconds")
			}
			var err error
			interval, err = strconv.Atoi(args[i+1])
			if err != nil {
				return fmt.Errorf("invalid timeout value: %s", args[i+1])
			}
			i += 2
		} else {
			break
		}
	}

	// Parse required arguments
	if i >= len(args) {
		return fmt.Errorf("usage: bsubio wait [-v] [-t <seconds>] <jobid>")
	}

	jobID = args[i]

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Poll for job completion
	if verbose {
		fmt.Printf("Waiting for job %s to complete (polling every %d seconds)...\n", jobID, interval)
	}

	for {
		resp, err := client.GetJobWithResponse(ctx, bsubio.JobId(jobID))
		if err != nil {
			return fmt.Errorf("failed to get job status: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("failed to get job status: HTTP %d", resp.StatusCode())
		}

		if resp.JSON200 == nil || resp.JSON200.Data == nil {
			return fmt.Errorf("unexpected response format")
		}

		job := resp.JSON200.Data

		if verbose && job.Status != nil {
			fmt.Printf("Status: %s\n", *job.Status)
		}

		// Check if job is in a terminal state
		if job.Status != nil {
			switch *job.Status {
			case "finished":
				fmt.Printf("Job completed successfully\n")
				return nil
			case "failed":
				if job.ErrorMessage != nil {
					return fmt.Errorf("job failed: %s", *job.ErrorMessage)
				}
				return fmt.Errorf("job failed")
			}
		}

		// Wait before polling again
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
