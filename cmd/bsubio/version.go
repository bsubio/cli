package main

import (
	"fmt"
)

func runVersion(args []string) error {
	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Get API server version
	resp, err := client.GetVersionWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get API version: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get API version: HTTP %d", resp.StatusCode())
	}

	fmt.Printf("CLI Version:    %s\n", version)

	if resp.JSON200 != nil && resp.JSON200.Version != nil {
		fmt.Printf("Server Version: %s\n", *resp.JSON200.Version)
	}

	return nil
}
