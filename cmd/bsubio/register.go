package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://app.bsub.io"
)

// runRegister implements the register command using GitHub Device Flow
func runRegister(args []string) error {
	fs := flag.NewFlagSet("register", flag.ContinueOnError)

	// Define flags
	verbose := fs.Bool("verbose", false, "Verbose output")
	debug := fs.Bool("debug", false, "Debug output")
	baseURLFlag := fs.String("base-url", "", "Base URL override (default: https://app.bsub.io)")

	// Custom usage function
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: bsubio register [options]\n\n")
		fmt.Fprintf(fs.Output(), "Register with bsub.io using GitHub authentication\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Determine base URL (priority: flag > environment > default)
	baseURL := *baseURLFlag
	if baseURL == "" {
		baseURL = os.Getenv("BSUBIO_BASE_URL")
		if baseURL == "" {
			baseURL = defaultBaseURL
		}
	}

	if *verbose || *debug {
		fmt.Printf("Using base URL: %s\n", baseURL)
		fmt.Printf("Hostname: %s\n", hostname)
	}

	fmt.Println("Registering with bsub.io using GitHub authentication...")
	fmt.Println()

	// Step 1: Request device code
	if *verbose || *debug {
		fmt.Printf("Requesting device code from %s/v1/auth/device/code\n", baseURL)
	}
	deviceCode, userCode, verificationURI, expiresIn, interval, err := requestDeviceCode(baseURL, hostname, *debug)
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}

	if *debug {
		fmt.Printf("Device code received (expires in %d seconds, poll interval: %d seconds)\n", expiresIn, interval)
	}

	// Step 2: Display code and prompt user
	fmt.Printf("! First copy your one-time code: %s\n", userCode)
	fmt.Printf("Press Enter to open %s in your browser...", verificationURI)

	// Wait for user to press Enter
	fmt.Scanln()

	// Open browser
	if *verbose || *debug {
		fmt.Printf("\nOpening browser to: %s\n", verificationURI)
	}
	if err := openBrowser(verificationURI); err != nil {
		fmt.Printf("\nCould not open browser automatically. Please visit the URL above manually.\n")
		if *debug {
			fmt.Printf("Browser error: %v\n", err)
		}
	}

	fmt.Println("Waiting for authorization...")

	// Step 3: Poll for authorization
	apiKey, userInfo, err := pollForAuthorization(baseURL, deviceCode, userCode, interval, expiresIn, *verbose, *debug)
	if err != nil {
		return fmt.Errorf("\nauthorization failed: %w", err)
	}

	fmt.Println("✓ Authentication complete.")

	// Step 4: Save configuration
	config := &Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✓ Logged in as %s %s (%s)\n", userInfo.FirstName, userInfo.LastName, userInfo.Email)

	return nil
}

// openBrowser opens the specified URL in the user's default browser
func openBrowser(urlStr string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid verification URL: %w", err)
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("verification URL must use HTTPS")
	}

	if !strings.HasSuffix(parsed.Host, ".bsub.io") && parsed.Host != "bsub.io" {
		return fmt.Errorf("unexpected verification host: %s", parsed.Host)
	}

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{urlStr}
	case "linux":
		cmd = "xdg-open"
		args = []string{urlStr}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", urlStr}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}

// newHTTPClient creates a secure HTTP client with TLS 1.2+ and timeout
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

// requestDeviceCode initiates the device flow
func requestDeviceCode(baseURL, hostname string, debug bool) (deviceCode, userCode, verificationURI string, expiresIn, interval int, err error) {
	endpoint := fmt.Sprintf("%s/v1/auth/device/code", baseURL)

	reqBody := map[string]string{
		"hostname": hostname,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	if debug {
		fmt.Printf("POST %s\n", endpoint)
		fmt.Printf("Request body: %s\n", string(jsonData))
	}

	client := newHTTPClient()

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", 0, 0, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && debug {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if debug {
		fmt.Printf("Response status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", "", 0, 0, fmt.Errorf("failed to read response: %w", err)
		}
		return "", "", "", 0, 0, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", "", 0, 0, err
	}

	return response.DeviceCode, response.UserCode, response.VerificationURI, response.ExpiresIn, response.Interval, nil
}

// pollForAuthorization polls the server until authorization is granted
func pollForAuthorization(baseURL, deviceCode, userCode string, interval, expiresIn int, verbose, debug bool) (apiKey string, userInfo *UserInfo, err error) {
	endpoint := fmt.Sprintf("%s/v1/auth/device/token", baseURL)
	pollInterval := time.Duration(interval) * time.Second
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)

	client := newHTTPClient()

	if debug {
		fmt.Printf("Polling %s every %d seconds until %s\n", endpoint, interval, deadline.Format(time.RFC3339))
	}

	for {
		if time.Now().After(deadline) {
			return "", nil, fmt.Errorf("authorization timeout - code expired")
		}

		reqBody := map[string]string{
			"device_code": deviceCode,
			"user_code":   userCode,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return "", nil, err
		}

		if debug {
			fmt.Printf("\nPOST %s\n", endpoint)
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "", nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if closeErr := resp.Body.Close(); closeErr != nil && debug {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
		if err != nil {
			return "", nil, err
		}

		if debug {
			fmt.Printf("Response status: %d\n", resp.StatusCode)
		}

		// Handle different response status codes
		switch resp.StatusCode {
		case http.StatusOK:
			// Authorization successful
			var response struct {
				Status string `json:"status"`
				APIKey string `json:"api_key"`
				User   struct {
					Email     string `json:"email"`
					FirstName string `json:"first_name"`
					LastName  string `json:"last_name"`
				} `json:"user"`
			}

			if err := json.Unmarshal(body, &response); err != nil {
				return "", nil, fmt.Errorf("failed to parse success response: %w", err)
			}

			if response.Status == "authorized" && response.APIKey != "" {
				userInfo := &UserInfo{
					Email:     response.User.Email,
					FirstName: response.User.FirstName,
					LastName:  response.User.LastName,
				}
				if verbose || debug {
					fmt.Printf("\nAuthorization successful for %s\n", userInfo.Email)
				}
				return response.APIKey, userInfo, nil
			}

		case http.StatusAccepted:
			// Still pending, continue polling
			fmt.Print(".")

		case http.StatusTooManyRequests:
			// Slow down
			var response struct {
				Interval int `json:"interval"`
			}
			if err := json.Unmarshal(body, &response); err == nil && response.Interval > 0 {
				if debug {
					fmt.Printf("\nRate limited, increasing poll interval to %d seconds\n", response.Interval)
				}
				pollInterval = time.Duration(response.Interval) * time.Second
			}
			fmt.Print(".")

		case http.StatusGone:
			// Code expired
			return "", nil, fmt.Errorf("authorization code expired")

		case http.StatusForbidden:
			// User denied access
			return "", nil, fmt.Errorf("authorization denied by user")

		default:
			// Other error
			return "", nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
		}

		// Wait before next poll
		time.Sleep(pollInterval)
	}
}

// UserInfo contains user information returned from registration
type UserInfo struct {
	Email     string
	FirstName string
	LastName  string
}
