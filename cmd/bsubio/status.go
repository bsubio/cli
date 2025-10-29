package main

import (
	"fmt"

	"github.com/bsubio/bsubio-go"
)

func runStatus(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bsubio status <jobid>")
	}

	jobID := args[0]

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get job details
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

	// Display job details
	fmt.Println("Job Details:")
	fmt.Println("--------------------------------------------------------------------------------")

	if job.Id != nil {
		fmt.Printf("ID:          %s\n", *job.Id)
	}

	if job.Type != nil {
		fmt.Printf("Type:        %s\n", *job.Type)
	}

	if job.Status != nil {
		fmt.Printf("Status:      %s\n", *job.Status)
	}

	if job.DataSize != nil {
		fmt.Printf("Data Size:   %d bytes\n", *job.DataSize)
	}

	if job.CreatedAt != nil {
		fmt.Printf("Created At:  %s\n", job.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if job.ClaimedAt != nil {
		fmt.Printf("Claimed At:  %s\n", job.ClaimedAt.Format("2006-01-02 15:04:05"))
	}

	if job.FinishedAt != nil {
		fmt.Printf("Finished At: %s\n", job.FinishedAt.Format("2006-01-02 15:04:05"))
	}

	if job.ClaimedBy != nil {
		fmt.Printf("Claimed By:  %s\n", *job.ClaimedBy)
	}

	if job.ErrorMessage != nil && *job.ErrorMessage != "" {
		fmt.Printf("\nError:       %s\n", *job.ErrorMessage)
	}

	return nil
}
