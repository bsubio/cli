package main

import (
	"fmt"

	"github.com/bsubio/bsubio-go"
)

func runRm(args []string) error {
	var (
		deleteAll bool
		jobID     string
	)

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-a" || arg == "--all" {
			deleteAll = true
			i++
		} else {
			break
		}
	}

	// If not deleting all, require job ID
	if !deleteAll {
		if i >= len(args) {
			return fmt.Errorf("usage: bsubio rm [-a|--all] <jobid>")
		}
		jobID = args[i]
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	if deleteAll {
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
				fmt.Printf("Failed to delete job %s: %v\n", *job.Id, err)
				continue
			}

			if deleteResp.StatusCode() != 200 && deleteResp.StatusCode() != 204 {
				fmt.Printf("Failed to delete job %s: HTTP %d\n", *job.Id, deleteResp.StatusCode())
				continue
			}

			fmt.Printf("Deleted job: %s\n", *job.Id)
			deletedCount++
		}

		fmt.Printf("Deleted %d job(s)\n", deletedCount)
	} else {
		// Delete single job
		resp, err := client.DeleteJobWithResponse(ctx, bsubio.JobId(jobID))
		if err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			return fmt.Errorf("failed to delete job: HTTP %d", resp.StatusCode())
		}

		fmt.Printf("Job deleted: %s\n", jobID)
	}

	return nil
}
