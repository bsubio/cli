package main

import (
	"flag"
	"fmt"
	"os"
)

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
				return fmt.Errorf("input file not found: %s", inputFile)
			}
			return fmt.Errorf("failed to access input file: %w", err)
		}
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Submit jobs for all input files
	jobIDs := make([]string, 0, len(inputFiles))
	for _, inputFile := range inputFiles {
		fmt.Fprintf(os.Stderr, "Submitting job for %s...\n", inputFile)
		job, err := client.CreateAndSubmitJobFromFile(ctx, jobType, inputFile)
		if err != nil {
			return fmt.Errorf("failed to submit job for %s: %w", inputFile, err)
		}
		fmt.Fprintf(os.Stderr, "Job submitted: %s\n", *job.Id)
		jobIDs = append(jobIDs, *job.Id)
	}

	// If wait flag is set, wait for all jobs and get outputs
	if *wait {
		for i, jobID := range jobIDs {
			fmt.Fprintf(os.Stderr, "Waiting for job %s to complete...\n", jobID)
			finishedJob, err := client.WaitForJob(ctx, jobID)
			if err != nil {
				return fmt.Errorf("failed to wait for job %s: %w", jobID, err)
			}

			if finishedJob.Status != nil && *finishedJob.Status == "failed" {
				if finishedJob.ErrorMessage != nil {
					return fmt.Errorf("job %s failed: %s", jobID, *finishedJob.ErrorMessage)
				}
				return fmt.Errorf("job %s failed", jobID)
			}

			fmt.Fprintf(os.Stderr, "Job %s completed successfully\n", jobID)

			// Get output
			outputResp, err := client.GetJobOutput(ctx, jobID)
			if err != nil {
				return fmt.Errorf("failed to get job output for %s: %w", jobID, err)
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
