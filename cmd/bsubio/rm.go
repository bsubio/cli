package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bsubio/bsubio-go"
	"github.com/google/uuid"
)

func runRm(args []string) error {
	fs := flag.NewFlagSet("rm", flag.ContinueOnError)

	// Define flags
	deleteAll := fs.Bool("a", false, "Delete all jobs")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio rm [options] [jobid]\n\n")
		fmt.Fprintf(fs.Output(), "Delete a job or all jobs\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nArguments:\n")
		fmt.Fprintf(fs.Output(), "  jobid    Job ID to delete (not required with -a)\n")
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Get remaining arguments
	remainingArgs := fs.Args()
	var jobID string

	// If not deleting all, require job ID
	if !*deleteAll {
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

	if *deleteAll {
		// List all jobs and delete them
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
		deletedCount := 0

		for _, job := range jobs {
			deleteResp, err := client.DeleteJobWithResponse(ctx, *job.Id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete job %s: %v\n", job.Id.String(), err)
				continue
			}

			if deleteResp.StatusCode() != 200 && deleteResp.StatusCode() != 204 {
				fmt.Fprintf(os.Stderr, "Failed to delete job %s: HTTP %d\n", job.Id.String(), deleteResp.StatusCode())
				continue
			}

			fmt.Fprintf(os.Stderr, "Deleted job: %s\n", job.Id.String())
			deletedCount++
		}

		fmt.Fprintf(os.Stderr, "Deleted %d job(s)\n", deletedCount)
	} else {
		// Delete single job
		jobUUID, err := uuid.Parse(jobID)
		if err != nil {
			return fmt.Errorf("invalid job ID: %w", err)
		}
		resp, err := client.DeleteJobWithResponse(ctx, jobUUID)
		if err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			return fmt.Errorf("failed to delete job: HTTP %d", resp.StatusCode())
		}

		fmt.Fprintf(os.Stderr, "Job deleted: %s\n", jobID)
	}

	return nil
}
