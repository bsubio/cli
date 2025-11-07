package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func runBench(args []string) error {
	fs := flag.NewFlagSet("bench", flag.ContinueOnError)

	// Define flags
	jobType := fs.String("type", "pdf_extract", "Job type to use for benchmarking")
	dataDir := fs.String("dir", "tests/data", "Directory containing test files")
	pattern := fs.String("pattern", "*.pdf", "File pattern to match (e.g., *.pdf)")
	jsonOutput := fs.Bool("json", false, "Output results in JSON format")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio bench [options]\n\n")
		fmt.Fprintf(fs.Output(), "Benchmark job processing with test files\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Check if data directory exists
	if _, err := os.Stat(*dataDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("data directory not found: %s", *dataDir)
		}
		return fmt.Errorf("failed to access data directory: %w", err)
	}

	// Find test files
	testFiles, err := filepath.Glob(filepath.Join(*dataDir, *pattern))
	if err != nil {
		return fmt.Errorf("failed to find test files: %w", err)
	}

	if len(testFiles) == 0 {
		return fmt.Errorf("no test files found in %s matching %s", *dataDir, *pattern)
	}

	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := getContext()

	if !*jsonOutput {
		fmt.Printf("Benchmarking %d file(s) with job type: %s\n", len(testFiles), *jobType)
		fmt.Println("================================================================================")
	}

	type result struct {
		File     string `json:"file"`
		Size     int64  `json:"size"`
		JobID    string `json:"job_id,omitempty"`
		SubmitMs int64  `json:"submit_ms"`
		TotalMs  int64  `json:"total_ms"`
		Status   string `json:"status"`
		Error    string `json:"error,omitempty"`
	}

	results := make([]result, 0, len(testFiles))

	// Process each test file
	for _, testFile := range testFiles {
		fileInfo, err := os.Stat(testFile)
		if err != nil {
			results = append(results, result{
				File:   filepath.Base(testFile),
				Status: "error",
				Error:  err.Error(),
			})
			continue
		}

		if !*jsonOutput {
			fmt.Printf("\nProcessing: %s (%d bytes)\n", filepath.Base(testFile), fileInfo.Size())
		}

		// Time submission
		submitStart := time.Now()
		job, err := client.CreateAndSubmitJobFromFile(ctx, *jobType, testFile)
		submitDuration := time.Since(submitStart)

		if err != nil {
			results = append(results, result{
				File:   filepath.Base(testFile),
				Size:   fileInfo.Size(),
				Status: "submit_failed",
				Error:  err.Error(),
			})
			if !*jsonOutput {
				fmt.Printf("  Submit failed: %v\n", err)
			}
			continue
		}

		if !*jsonOutput {
			fmt.Printf("  Job ID: %s\n", *job.Id)
			fmt.Printf("  Submit time: %dms\n", submitDuration.Milliseconds())
		}

		// Wait for completion
		finishedJob, err := client.WaitForJob(ctx, *job.Id)
		totalDuration := time.Since(submitStart)

		if err != nil {
			results = append(results, result{
				File:     filepath.Base(testFile),
				Size:     fileInfo.Size(),
				JobID:    *job.Id,
				SubmitMs: submitDuration.Milliseconds(),
				Status:   "wait_failed",
				Error:    err.Error(),
			})
			if !*jsonOutput {
				fmt.Printf("  Wait failed: %v\n", err)
			}
			continue
		}

		jobStatus := "unknown"
		if finishedJob.Status != nil {
			jobStatus = string(*finishedJob.Status)
		}

		errorMsg := ""
		if jobStatus == "failed" && finishedJob.ErrorMessage != nil {
			errorMsg = *finishedJob.ErrorMessage
		}

		results = append(results, result{
			File:     filepath.Base(testFile),
			Size:     fileInfo.Size(),
			JobID:    *job.Id,
			SubmitMs: submitDuration.Milliseconds(),
			TotalMs:  totalDuration.Milliseconds(),
			Status:   jobStatus,
			Error:    errorMsg,
		})

		if !*jsonOutput {
			fmt.Printf("  Status: %s\n", jobStatus)
			fmt.Printf("  Total time: %dms\n", totalDuration.Milliseconds())

			if errorMsg != "" {
				fmt.Printf("  Error: %s\n", errorMsg)
			}
		}
	}

	// Output results
	if *jsonOutput {
		// JSON output
		type jsonOutput struct {
			JobType      string   `json:"job_type"`
			TotalFiles   int      `json:"total_files"`
			Successful   int      `json:"successful"`
			AvgSubmitMs  int64    `json:"avg_submit_ms"`
			AvgTotalMs   int64    `json:"avg_total_ms"`
			Results      []result `json:"results"`
		}

		var totalSubmitMs int64
		var totalProcessMs int64
		successCount := 0

		for _, r := range results {
			totalSubmitMs += r.SubmitMs
			totalProcessMs += r.TotalMs
			if r.Status == "finished" || r.Status == "completed" {
				successCount++
			}
		}

		avgSubmitMs := int64(0)
		avgTotalMs := int64(0)
		if len(results) > 0 {
			avgSubmitMs = totalSubmitMs / int64(len(results))
			avgTotalMs = totalProcessMs / int64(len(results))
		}

		output := jsonOutput{
			JobType:     *jobType,
			TotalFiles:  len(results),
			Successful:  successCount,
			AvgSubmitMs: avgSubmitMs,
			AvgTotalMs:  avgTotalMs,
			Results:     results,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output
		fmt.Println("\n================================================================================")
		fmt.Println("SUMMARY")
		fmt.Println("================================================================================")
		fmt.Printf("%-30s %10s %10s %10s %s\n", "File", "Size", "Submit", "Total", "Status")
		fmt.Println("--------------------------------------------------------------------------------")

		var totalSubmitMs int64
		var totalProcessMs int64
		successCount := 0

		for _, r := range results {
			if r.Error != "" {
				fmt.Printf("%-30s %10s %10s %10s %s: %s\n",
					truncate(r.File, 30),
					formatBytes(r.Size),
					"-",
					"-",
					r.Status,
					r.Error)
			} else {
				fmt.Printf("%-30s %10s %10dms %10dms %s\n",
					truncate(r.File, 30),
					formatBytes(r.Size),
					r.SubmitMs,
					r.TotalMs,
					r.Status)

				totalSubmitMs += r.SubmitMs
				totalProcessMs += r.TotalMs
				if r.Status == "finished" || r.Status == "completed" {
					successCount++
				}
			}
		}

		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("Successful: %d/%d\n", successCount, len(results))
		if len(results) > 0 {
			fmt.Printf("Avg Submit: %dms\n", totalSubmitMs/int64(len(results)))
			fmt.Printf("Avg Total:  %dms\n", totalProcessMs/int64(len(results)))
		}
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}
