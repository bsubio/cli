package main

import (
	"fmt"
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

	fmt.Printf("%-20s %-20s %-20s %s\n", "Worker Type", "MIME in", "MIME out", "Description")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, jobType := range types {
		workerType := ""
		if jobType.Type != nil {
			workerType = *jobType.Type
		}

		description := ""
		if jobType.Description != nil {
			description = *jobType.Description
		}

		// Get MIME out value
		mimeOut := ""
		if jobType.Output != nil && jobType.Output.MimeOut != nil {
			mimeOut = *jobType.Output.MimeOut
		}

		// Print a row for each MIME in value
		if jobType.Input != nil && jobType.Input.MimeIn != nil {
			for _, mimeIn := range *jobType.Input.MimeIn {
				fmt.Printf("%-20s %-20s %-20s %s\n", workerType, mimeIn, mimeOut, description)
			}
			continue
		}

		// If no MIME in values, print one row
		fmt.Printf("%-20s %-20s %-20s %s\n", workerType, "", mimeOut, description)
	}

	return nil
}
