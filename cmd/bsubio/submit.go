package main

import (
	"fmt"
	"os"
)

func runSubmit(args []string) error {
	var (
		wait       bool
		outputFile string
		inputFile  string
		jobType    string
	)

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-w" {
			wait = true
			i++
		} else if arg == "-o" {
			if i+1 >= len(args) {
				return fmt.Errorf("-o flag requires a file path")
			}
			outputFile = args[i+1]
			i += 2
		} else {
			break
		}
	}

	// Parse required arguments
	if i+2 > len(args) {
		return fmt.Errorf("usage: bsubio submit [-o <file>] [-w] <input_file> <type>")
	}

	inputFile = args[i]
	jobType = args[i+1]

	// Validate that output file is only used with wait
	if outputFile != "" && !wait {
		return fmt.Errorf("-o flag requires -w flag")
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputFile)
		}
		return fmt.Errorf("failed to access input file: %w", err)
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Submit job
	fmt.Printf("Submitting job...\n")
	job, err := client.CreateAndSubmitJobFromFile(ctx, jobType, inputFile)
	if err != nil {
		return fmt.Errorf("failed to submit job: %w", err)
	}

	fmt.Printf("Job submitted: %s\n", *job.Id)

	// If wait flag is set, wait for completion and get output
	if wait {
		fmt.Printf("Waiting for job to complete...\n")
		finishedJob, err := client.WaitForJob(ctx, *job.Id)
		if err != nil {
			return fmt.Errorf("failed to wait for job: %w", err)
		}

		if finishedJob.Status != nil && *finishedJob.Status == "failed" {
			if finishedJob.ErrorMessage != nil {
				return fmt.Errorf("job failed: %s", *finishedJob.ErrorMessage)
			}
			return fmt.Errorf("job failed")
		}

		fmt.Printf("Job completed successfully\n")

		// Get output
		outputResp, err := client.GetJobOutput(ctx, *job.Id)
		if err != nil {
			return fmt.Errorf("failed to get job output: %w", err)
		}
		defer outputResp.Body.Close()

		// Write output to file or stdout
		if outputFile != "" {
			file, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer file.Close()

			if _, err := file.ReadFrom(outputResp.Body); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Printf("Output saved to %s\n", outputFile)
		} else {
			if _, err := os.Stdout.ReadFrom(outputResp.Body); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}
	}

	return nil
}
