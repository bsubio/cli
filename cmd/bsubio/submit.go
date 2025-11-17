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
		fmt.Fprintf(fs.Output(), "Usage: bsubio submit [options] <type> <input_file> [<input_file2> ...]\n\n")
		fmt.Fprintf(fs.Output(), "Submit a job for processing\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nArguments:\n")
		fmt.Fprintf(fs.Output(), "  type          Job type (e.g., pdf_extract)\n")
		fmt.Fprintf(fs.Output(), "  input_file    Path to the input file(s)\n")
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

	// Validate that output file is only used with wait
	if *outputFile != "" && !*wait {
		return fmt.Errorf("-o flag requires -w flag")
	}

	// Validate that output file is not used with multiple files
	if *outputFile != "" && len(inputFiles) > 1 {
		return fmt.Errorf("-o flag cannot be used with multiple input files")
	}

	// Check if all input files exist
	for _, inputFile := range inputFiles {
		if _, err := os.Stat(inputFile); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("input file not found: %s", inputFile)
			}
			return fmt.Errorf("failed to access input file %s: %w", inputFile, err)
		}
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	// Submit all jobs and collect their IDs
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

	// If wait flag is set, wait for all jobs and get their outputs
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
					outputResp.Body.Close()
					return fmt.Errorf("failed to create output file: %w", err)
				}

				if _, err := file.ReadFrom(outputResp.Body); err != nil {
					file.Close()
					outputResp.Body.Close()
					return fmt.Errorf("failed to write output file: %w", err)
				}

				file.Close()
				outputResp.Body.Close()
				fmt.Fprintf(os.Stderr, "Output saved to %s\n", *outputFile)
			} else {
				// For multiple files, add separator between outputs
				if len(inputFiles) > 1 && i > 0 {
					fmt.Fprintf(os.Stdout, "\n")
				}
				if _, err := os.Stdout.ReadFrom(outputResp.Body); err != nil {
					outputResp.Body.Close()
					return fmt.Errorf("failed to write output for job %s: %w", jobID, err)
				}
				outputResp.Body.Close()
			}
		}
	}

	return nil
}
