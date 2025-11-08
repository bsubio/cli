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
	// Check if this is a diff subcommand
	if len(args) > 0 && args[0] == "diff" {
		return runBenchDiff(args[1:])
	}

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

	results := make([]benchResult, 0, len(testFiles))

	// Process each test file
	for _, testFile := range testFiles {
		fileInfo, err := os.Stat(testFile)
		if err != nil {
			results = append(results, benchResult{
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
			results = append(results, benchResult{
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
			results = append(results, benchResult{
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

		results = append(results, benchResult{
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

		output := benchOutput{
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

type benchResult struct {
	File     string `json:"file"`
	Size     int64  `json:"size"`
	JobID    string `json:"job_id,omitempty"`
	SubmitMs int64  `json:"submit_ms"`
	TotalMs  int64  `json:"total_ms"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

type benchOutput struct {
	JobType     string        `json:"job_type"`
	TotalFiles  int           `json:"total_files"`
	Successful  int           `json:"successful"`
	AvgSubmitMs int64         `json:"avg_submit_ms"`
	AvgTotalMs  int64         `json:"avg_total_ms"`
	Results     []benchResult `json:"results"`
}

func runBenchDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: bsubio bench diff <file1.json> <file2.json>")
	}

	file1Path := args[0]
	file2Path := args[1]

	// Read first benchmark file
	bench1, err := readBenchmarkFile(file1Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1Path, err)
	}

	// Read second benchmark file
	bench2, err := readBenchmarkFile(file2Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2Path, err)
	}

	// Create a map for quick lookup
	bench2Map := make(map[string]benchResult)
	for _, r := range bench2.Results {
		bench2Map[r.File] = r
	}

	// Print header
	fmt.Printf("Comparing: %s vs %s\n", filepath.Base(file1Path), filepath.Base(file2Path))
	fmt.Printf("Job types: %s vs %s\n\n", bench1.JobType, bench2.JobType)

	fmt.Printf("%-40s %20s %20s %12s\n", "Filename",
		bench1.JobType,
		bench2.JobType,
		"Diff (%)")
	fmt.Println("----------------------------------------------------------------------------------------------------")

	// Compare results
	for _, r1 := range bench1.Results {
		r2, found := bench2Map[r1.File]

		if !found {
			fmt.Printf("%-40s %20s %20s %12s\n",
				r1.File,
				fmt.Sprintf("%dms", r1.TotalMs),
				"-",
				"-")
			continue
		}

		// Skip failed jobs
		if r1.Status != "finished" && r1.Status != "completed" {
			continue
		}
		if r2.Status != "finished" && r2.Status != "completed" {
			continue
		}

		// Calculate percentage difference
		var diffPct float64
		if r1.TotalMs > 0 {
			diffPct = ((float64(r2.TotalMs) - float64(r1.TotalMs)) / float64(r1.TotalMs)) * 100
		}

		diffStr := fmt.Sprintf("%.1f%%", diffPct)
		if diffPct > 0 {
			diffStr = "+" + diffStr
		}

		fmt.Printf("%-40s %20dms %20dms %12s\n",
			r1.File,
			r1.TotalMs,
			r2.TotalMs,
			diffStr)
	}

	// Print summary
	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Printf("%-40s %20dms %20dms",
		"Average",
		bench1.AvgTotalMs,
		bench2.AvgTotalMs)

	if bench1.AvgTotalMs > 0 {
		avgDiffPct := ((float64(bench2.AvgTotalMs) - float64(bench1.AvgTotalMs)) / float64(bench1.AvgTotalMs)) * 100
		diffStr := fmt.Sprintf("%.1f%%", avgDiffPct)
		if avgDiffPct > 0 {
			diffStr = "+" + diffStr
		}
		fmt.Printf(" %12s\n", diffStr)
	} else {
		fmt.Printf(" %12s\n", "-")
	}

	return nil
}

func readBenchmarkFile(path string) (*benchOutput, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var output benchOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func truncateJobType(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// Remove common prefixes/suffixes to shorten
	s = filepath.Base(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}
