package main

import (
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

	fmt.Printf("Benchmarking %d file(s) with job type: %s\n", len(testFiles), *jobType)
	fmt.Println("================================================================================")

	type result struct {
		file     string
		size     int64
		jobID    string
		submitMs int64
		totalMs  int64
		status   string
		err      error
	}

	results := make([]result, 0, len(testFiles))

	// Process each test file
	for _, testFile := range testFiles {
		fileInfo, err := os.Stat(testFile)
		if err != nil {
			results = append(results, result{
				file:   filepath.Base(testFile),
				status: "error",
				err:    err,
			})
			continue
		}

		fmt.Printf("\nProcessing: %s (%d bytes)\n", filepath.Base(testFile), fileInfo.Size())

		// Time submission
		submitStart := time.Now()
		job, err := client.CreateAndSubmitJobFromFile(ctx, *jobType, testFile)
		submitDuration := time.Since(submitStart)

		if err != nil {
			results = append(results, result{
				file:   filepath.Base(testFile),
				size:   fileInfo.Size(),
				status: "submit_failed",
				err:    err,
			})
			fmt.Printf("  Submit failed: %v\n", err)
			continue
		}

		fmt.Printf("  Job ID: %s\n", *job.Id)
		fmt.Printf("  Submit time: %dms\n", submitDuration.Milliseconds())

		// Wait for completion
		finishedJob, err := client.WaitForJob(ctx, *job.Id)
		totalDuration := time.Since(submitStart)

		if err != nil {
			results = append(results, result{
				file:     filepath.Base(testFile),
				size:     fileInfo.Size(),
				jobID:    *job.Id,
				submitMs: submitDuration.Milliseconds(),
				status:   "wait_failed",
				err:      err,
			})
			fmt.Printf("  Wait failed: %v\n", err)
			continue
		}

		jobStatus := "unknown"
		if finishedJob.Status != nil {
			jobStatus = string(*finishedJob.Status)
		}

		results = append(results, result{
			file:     filepath.Base(testFile),
			size:     fileInfo.Size(),
			jobID:    *job.Id,
			submitMs: submitDuration.Milliseconds(),
			totalMs:  totalDuration.Milliseconds(),
			status:   jobStatus,
		})

		fmt.Printf("  Status: %s\n", jobStatus)
		fmt.Printf("  Total time: %dms\n", totalDuration.Milliseconds())

		if jobStatus == "failed" && finishedJob.ErrorMessage != nil {
			fmt.Printf("  Error: %s\n", *finishedJob.ErrorMessage)
		}
	}

	// Print summary
	fmt.Println("\n================================================================================")
	fmt.Println("SUMMARY")
	fmt.Println("================================================================================")
	fmt.Printf("%-30s %10s %10s %10s %s\n", "File", "Size", "Submit", "Total", "Status")
	fmt.Println("--------------------------------------------------------------------------------")

	var totalSubmitMs int64
	var totalProcessMs int64
	successCount := 0

	for _, r := range results {
		if r.err != nil {
			fmt.Printf("%-30s %10s %10s %10s %s: %v\n",
				truncate(r.file, 30),
				formatBytes(r.size),
				"-",
				"-",
				r.status,
				r.err)
		} else {
			fmt.Printf("%-30s %10s %10dms %10dms %s\n",
				truncate(r.file, 30),
				formatBytes(r.size),
				r.submitMs,
				r.totalMs,
				r.status)

			totalSubmitMs += r.submitMs
			totalProcessMs += r.totalMs
			if r.status == "finished" || r.status == "completed" {
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
