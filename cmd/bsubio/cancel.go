package main

import (
	"fmt"

	"github.com/bsubio/bsubio-go"
)

func runCancel(args []string) error {
	var (
		cancelAll bool
		jobID     string
	)

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-a" || arg == "--all" {
			cancelAll = true
			i++
		} else {
			break
		}
	}

	// If not canceling all, require job ID
	if !cancelAll {
		if i >= len(args) {
			return fmt.Errorf("usage: bsubio cancel [-a|--all] <jobid>")
		}
		jobID = args[i]
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	if cancelAll {
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
					fmt.Printf("Failed to cancel job %s: %v\n", *job.Id, err)
					continue
				}

				if cancelResp.StatusCode() != 200 {
					fmt.Printf("Failed to cancel job %s: HTTP %d\n", *job.Id, cancelResp.StatusCode())
					continue
				}

				fmt.Printf("Canceled job: %s\n", *job.Id)
				canceledCount++
			}
		}

		fmt.Printf("Canceled %d job(s)\n", canceledCount)
	} else {
		// Cancel single job
		resp, err := client.CancelJobWithResponse(ctx, bsubio.JobId(jobID))
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
