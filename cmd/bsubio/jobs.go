package main

import (
	"fmt"
	"strconv"

	"github.com/bsubio/bsubio-go"
)

func runJobs(args []string) error {
	var (
		status string
		limit  = 20 // default limit
	)

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		switch arg {
		case "--status":
			if i+1 >= len(args) {
				return fmt.Errorf("--status flag requires a value")
			}
			status = args[i+1]
			i += 2
		case "--limit":
			if i+1 >= len(args) {
				return fmt.Errorf("--limit flag requires a value")
			}
			var err error
			limit, err = strconv.Atoi(args[i+1])
			if err != nil {
				return fmt.Errorf("invalid limit value: %s", args[i+1])
			}
			i += 2
		default:
			return fmt.Errorf("unknown flag: %s", arg)
		}
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Build parameters
	params := &bsubio.ListJobsParams{
		Limit: &limit,
	}

	if status != "" {
		statusParam := bsubio.ListJobsParamsStatus(status)
		params.Status = &statusParam
	}

	// List jobs
	resp, err := client.ListJobsWithResponse(ctx, params)
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

	// Display jobs
	if len(jobs) == 0 {
		fmt.Println("No jobs found")
		return nil
	}

	// Find the longest type string for proper alignment
	maxTypeLen := len("TYPE")
	for _, job := range jobs {
		if job.Type != nil && len(*job.Type) > maxTypeLen {
			maxTypeLen = len(*job.Type)
		}
	}

	// Print header with dynamic TYPE column width
	fmt.Printf("%-40s %-*s %-15s %s\n", "JOB ID", maxTypeLen, "TYPE", "STATUS", "CREATED AT")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, job := range jobs {
		jobID := ""
		if job.Id != nil {
			jobID = job.Id.String()
		}

		jobType := ""
		if job.Type != nil {
			jobType = *job.Type
		}

		status := ""
		if job.Status != nil {
			status = string(*job.Status)
		}

		createdAt := ""
		if job.CreatedAt != nil {
			createdAt = job.CreatedAt.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-40s %-*s %-15s %s\n", jobID, maxTypeLen, jobType, status, createdAt)
	}

	return nil
}
