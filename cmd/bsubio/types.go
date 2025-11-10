package main

import (
	"fmt"
	"strings"
)

func runTypes(args []string) error {
	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get available job types
	resp, err := client.GetTypesWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get job types: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get job types: HTTP %d", resp.StatusCode())
	}

	if resp.JSON200 == nil || resp.JSON200.Types == nil {
		return fmt.Errorf("unexpected response format")
	}

	types := *resp.JSON200.Types

	if len(types) == 0 {
		fmt.Println("No job types available")
		return nil
	}

	fmt.Printf("%-20s %-30s %s\n", "Worker Type", "MIME", "Description")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, jobType := range types {
		workerType := ""
		if jobType.Type != nil {
			workerType = *jobType.Type
		}

		mime := ""
		if jobType.Mime != nil && len(*jobType.Mime) > 0 {
			mime = strings.Join(*jobType.Mime, ", ")
		}

		description := ""
		if jobType.Description != nil {
			description = *jobType.Description
		}

		fmt.Printf("%-20s %-30s %s\n", workerType, mime, description)
	}

	return nil
}
