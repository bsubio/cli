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

		// Get MIME in values
		var mimeInValues []string
		if jobType.Input != nil && jobType.Input.MimeIn != nil {
			mimeInValues = *jobType.Input.MimeIn
		}

		// Get MIME out value
		mimeOut := ""
		if jobType.Output != nil && jobType.Output.MimeOut != nil {
			mimeOut = *jobType.Output.MimeOut
		}

		// If no MIME in values, print one row
		if len(mimeInValues) == 0 {
			fmt.Printf("%-20s %-20s %-20s %s\n", workerType, "", mimeOut, description)
			continue
		}

		// Print a row for each MIME in value
		for _, mimeIn := range mimeInValues {
			fmt.Printf("%-20s %-20s %-20s %s\n", workerType, mimeIn, mimeOut, description)
		}
	}

	return nil
}
