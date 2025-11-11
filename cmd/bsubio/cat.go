package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func runCat(args []string) error {
	fs := flag.NewFlagSet("cat", flag.ContinueOnError)
	wait := fs.Bool("wait", false, "Wait for job to complete before showing output")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio cat [options] <jobid>\n\n")
		fmt.Fprintf(fs.Output(), "Print job output (stdout)\n\n")
		fmt.Fprintf(fs.Output(), "Arguments:\n")
		fmt.Fprintf(fs.Output(), "  jobid    Job ID\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
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

	// Check job status first
	statusResp, err := client.GetJobWithResponse(ctx, jobUUID)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	if statusResp.StatusCode() != 200 {
		return fmt.Errorf("failed to get job status: HTTP %d", statusResp.StatusCode())
	}

	if statusResp.JSON200 == nil || statusResp.JSON200.Data == nil {
		return fmt.Errorf("job data is missing")
	}

	job := statusResp.JSON200.Data
	if job.Status == nil {
		return fmt.Errorf("job status is missing")
	}

	// If job is not completed and wait is not set, return helpful error
	if *job.Status != "finished" && *job.Status != "failed" {
		if *wait {
			fmt.Fprintf(os.Stderr, "Job is %s, waiting for completion...\n", *job.Status)
			finishedJob, err := client.WaitForJob(ctx, jobUUID)
			if err != nil {
				return fmt.Errorf("failed to wait for job: %w", err)
			}

			if finishedJob.Status != nil && *finishedJob.Status == "failed" {
				if finishedJob.ErrorMessage != nil {
					return fmt.Errorf("job failed: %s", *finishedJob.ErrorMessage)
				}
				return fmt.Errorf("job failed")
			}
		} else {
			return fmt.Errorf("job is not complete (status: %s). Use 'bsubio wait %s' first or use -wait flag", *job.Status, jobID)
		}
	}

	// Get job output
	resp, err := client.GetJobOutput(ctx, jobUUID)
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
