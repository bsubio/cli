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

	fmt.Println("Available Job Types:")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, jobType := range types {
		name := ""
		if jobType.Name != nil {
			name = *jobType.Name
		}

		description := ""
		if jobType.Description != nil {
			description = *jobType.Description
		}

		fmt.Printf("%-20s %s\n", name, description)
	}

	return nil
}
