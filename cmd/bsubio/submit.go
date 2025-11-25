package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bsubio/bsubio-go"
)

func getAvailableTypes(client *bsubio.BsubClient) ([]string, error) {
	ctx := getContext()
	resp, err := client.GetTypesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get job types: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get job types: HTTP %d", resp.StatusCode())
	}

	if resp.JSON200 == nil || resp.JSON200.Types == nil {
		return nil, fmt.Errorf("unexpected response format")
	}

	types := make([]string, 0, len(*resp.JSON200.Types))
	for _, jobType := range *resp.JSON200.Types {
		if jobType.Type != nil {
			types = append(types, *jobType.Type)
		}
	}
	return types, nil
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	prev := make([]int, len(s2)+1)
	curr := make([]int, len(s2)+1)

	for j := 0; j <= len(s2); j++ {
		prev[j] = j
	}

	for i := 1; i <= len(s1); i++ {
		curr[0] = i
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			curr[j] = min(
				curr[j-1]+1,
				min(prev[j]+1, prev[j-1]+cost),
			)
		}
		prev, curr = curr, prev
	}
	return prev[len(s2)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findClosestType(jobType string, availableTypes []string) string {
	if len(availableTypes) == 0 {
		return ""
	}

	minDistance := levenshteinDistance(jobType, availableTypes[0])
	closest := availableTypes[0]

	for _, t := range availableTypes[1:] {
		distance := levenshteinDistance(jobType, t)
		if distance < minDistance {
			minDistance = distance
			closest = t
		}
	}
	return closest
}

func validateJobType(jobType string, availableTypes []string) error {
	for _, t := range availableTypes {
		if t == jobType {
			return nil
		}
	}

	suggestion := findClosestType(jobType, availableTypes)
	var typesMsg strings.Builder
	typesMsg.WriteString(fmt.Sprintf("invalid job type: %s\n\n", jobType))
	typesMsg.WriteString(fmt.Sprintf("Did you mean: %s\n\n", suggestion))
	typesMsg.WriteString("Available types:\n")
	for _, t := range availableTypes {
		typesMsg.WriteString(fmt.Sprintf("  - %s\n", t))
	}
	return fmt.Errorf("%s", typesMsg.String())
}

func runSubmit(args []string) error {
	fs := flag.NewFlagSet("submit", flag.ContinueOnError)

	// Define flags
	wait := fs.Bool("w", false, "Wait for job to complete")
	outputFile := fs.String("o", "", "Output file path (requires -w)")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio submit [options] <type> <input_file> [input_file2 ...]\n\n")
		fmt.Fprintf(fs.Output(), "Submit one or more jobs for processing\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nArguments:\n")
		fmt.Fprintf(fs.Output(), "  type          Job type (e.g., pdf_extract)\n")
		fmt.Fprintf(fs.Output(), "  input_file    Path to input file(s) (can specify multiple)\n")
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Get remaining arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) < 2 {
		fs.Usage()
		return fmt.Errorf("expected at least 2 arguments, got %d", len(remainingArgs))
	}

	jobType := remainingArgs[0]
	inputFiles := remainingArgs[1:]

	// Validate that output file is only used with wait and single file
	if *outputFile != "" && !*wait {
		return fmt.Errorf("-o flag requires -w flag")
	}
	if *outputFile != "" && len(inputFiles) > 1 {
		return fmt.Errorf("-o flag cannot be used with multiple input files")
	}

	// Check if all input files exist
	for _, inputFile := range inputFiles {
		if _, err := os.Stat(inputFile); err != nil {
			if os.IsNotExist(err) {
				return &ExitError{
					Message: fmt.Sprintf("input file not found: %s", inputFile),
					Code:    2,
				}
			}
			return fmt.Errorf("failed to access input file: %w", err)
		}
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	// Validate job type
	availableTypes, err := getAvailableTypes(client)
	if err != nil {
		return err
	}

	if err := validateJobType(jobType, availableTypes); err != nil {
		return err
	}

	ctx := getContext()

	// Submit jobs for all input files
	jobIDs := make([]bsubio.JobId, 0, len(inputFiles))
	for _, inputFile := range inputFiles {
		fmt.Fprintf(os.Stderr, "Submitting job for %s...\n", inputFile)
		job, err := client.CreateAndSubmitJobFromFile(ctx, jobType, inputFile)
		if err != nil {
			return fmt.Errorf("failed to submit job for %s: %w", inputFile, err)
		}
		fmt.Fprintf(os.Stderr, "Job submitted: %s\n", job.Id.String())
		jobIDs = append(jobIDs, *job.Id)
	}

	// If wait flag is set, wait for all jobs and get outputs
	if *wait {
		for i, jobID := range jobIDs {
			fmt.Fprintf(os.Stderr, "Waiting for job %s to complete...\n", jobID.String())
			finishedJob, err := client.WaitForJob(ctx, jobID)
			if err != nil {
				return fmt.Errorf("failed to wait for job %s: %w", jobID.String(), err)
			}

			if finishedJob.Status != nil && *finishedJob.Status == "failed" {
				if finishedJob.ErrorMessage != nil {
					return fmt.Errorf("job %s failed: %s", jobID.String(), *finishedJob.ErrorMessage)
				}
				return fmt.Errorf("job %s failed", jobID.String())
			}

			fmt.Fprintf(os.Stderr, "Job %s completed successfully\n", jobID.String())

			// Get output
			outputResp, err := client.GetJobOutput(ctx, jobID)
			if err != nil {
				return fmt.Errorf("failed to get job output for %s: %w", jobID.String(), err)
			}

			if outputResp.StatusCode != 200 {
				body, err := io.ReadAll(outputResp.Body)
				_ = outputResp.Body.Close()
				if err != nil {
					return fmt.Errorf("failed to get job output for %s: HTTP %d", jobID.String(), outputResp.StatusCode)
				}

				var errorResp struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
					return fmt.Errorf("failed to get job output for %s: %s", jobID.String(), errorResp.Error)
				}

				return fmt.Errorf("failed to get job output for %s: HTTP %d: %s", jobID.String(), outputResp.StatusCode, string(body))
			}

			// Write output to file or stdout
			if *outputFile != "" {
				file, err := os.Create(*outputFile)
				if err != nil {
					_ = outputResp.Body.Close()
					return fmt.Errorf("failed to create output file: %w", err)
				}

				if _, err := file.ReadFrom(outputResp.Body); err != nil {
					_ = file.Close()
					_ = outputResp.Body.Close()
					return fmt.Errorf("failed to write output file: %w", err)
				}

				_ = file.Close()
				_ = outputResp.Body.Close()
				fmt.Fprintf(os.Stderr, "Output saved to %s\n", *outputFile)
			} else {
				if _, err := os.Stdout.ReadFrom(outputResp.Body); err != nil {
					_ = outputResp.Body.Close()
					return fmt.Errorf("failed to write output: %w", err)
				}
				_ = outputResp.Body.Close()

				// Add separator between outputs if there are multiple files
				if len(inputFiles) > 1 && i < len(jobIDs)-1 {
					fmt.Println()
				}
			}
		}
	}

	return nil
}
