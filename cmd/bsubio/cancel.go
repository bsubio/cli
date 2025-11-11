package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
	"github.com/google/uuid"
)

func runCancel(args []string) error {
	fs := flag.NewFlagSet("cancel", flag.ContinueOnError)

	// Define flags
	cancelAll := fs.Bool("a", false, "Cancel all pending/claimed jobs")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio cancel [options] [jobid]\n\n")
		fmt.Fprintf(fs.Output(), "Cancel a job or all jobs\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nArguments:\n")
		fmt.Fprintf(fs.Output(), "  jobid    Job ID to cancel (not required with -a)\n")
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Get remaining arguments
	remainingArgs := fs.Args()
	var jobID string

	// If not canceling all, require job ID
	if !*cancelAll {
		if len(remainingArgs) != 1 {
			fs.Usage()
			return fmt.Errorf("expected 1 argument when not using -a flag")
		}
		jobID = remainingArgs[0]
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	if *cancelAll {
		// List all jobs and cancel pending/claimed ones
		resp, err := client.ListJobsWithResponse(ctx, &bsubio.ListJobsParams{})
		if err != nil {
			return fmt.Errorf("failed to list jobs: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("failed to list jobs: HTTP %d", resp.StatusCode())
		}

		if resp.JSON200 == nil || resp.JSON200.Data == nil || resp.JSON200.Data.Jobs == nil {
			return fmt.Errorf("unexpected response format")
		}

		jobs := *resp.JSON200.Data.Jobs
		canceledCount := 0

		for _, job := range jobs {
			if job.Status != nil && (*job.Status == "pending" || *job.Status == "claimed") {
				cancelResp, err := client.CancelJobWithResponse(ctx, *job.Id)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to cancel job %s: %v\n", *job.Id, err)
					continue
				}

				if cancelResp.StatusCode() != 200 {
					fmt.Fprintf(os.Stderr, "Failed to cancel job %s: HTTP %d\n", *job.Id, cancelResp.StatusCode())
					continue
				}

				fmt.Printf("Canceled job: %s\n", *job.Id)
				canceledCount++
			}
		}

		fmt.Printf("Canceled %d job(s)\n", canceledCount)
	} else {
		// Cancel single job
		jobUUID, err := uuid.Parse(jobID)
		if err != nil {
			return fmt.Errorf("invalid job ID: %w", err)
		}
		resp, err := client.CancelJobWithResponse(ctx, jobUUID)
		if err != nil {
			return fmt.Errorf("failed to cancel job: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("failed to cancel job: HTTP %d", resp.StatusCode())
		}

		fmt.Printf("Job canceled: %s\n", jobID)
	}

	return nil
}
