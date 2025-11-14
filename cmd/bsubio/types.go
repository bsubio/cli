package main

import (
	"fmt"
)

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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
		workerType := derefString(jobType.Type)
		description := derefString(jobType.Description)

		mimeOuts := []string{""}
		if jobType.Output != nil && jobType.Output.MimeOut != nil {
			mimeOuts = *jobType.Output.MimeOut
		}

		mimeIns := []string{""}
		if jobType.Input != nil && jobType.Input.MimeIn != nil {
			mimeIns = *jobType.Input.MimeIn
		}

		for _, mimeIn := range mimeIns {
			for _, mimeOut := range mimeOuts {
				fmt.Printf("%-20s %-20s %-20s %s\n", workerType, mimeIn, mimeOut, description)
			}
		}
	}

	return nil
}
