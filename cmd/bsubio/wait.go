package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bsubio/bsubio-go"
)

func runWait(args []string) error {
	fs := flag.NewFlagSet("wait", flag.ContinueOnError)

	// Define flags
	verbose := fs.Bool("v", false, "Verbose output")
	interval := fs.Int("t", 5, "Polling interval in seconds")

	// Custom usage function
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), "Usage: bsubio wait [options] <jobid>\n\n")
		_, _ = fmt.Fprintf(fs.Output(), "Wait for a job to complete\n\n")
		_, _ = fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		_, _ = fmt.Fprintf(fs.Output(), "\nArguments:\n")
		_, _ = fmt.Fprintf(fs.Output(), "  jobid    Job ID to wait for\n")
	}

	// Parse flags
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

	// Poll for job completion
	if *verbose {
		fmt.Printf("Waiting for job %s to complete (polling every %d seconds)...\n", jobID, *interval)
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

		if *verbose && job.Status != nil {
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
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
