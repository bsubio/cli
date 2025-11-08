package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	defaultBaseURL = "https://app.bsub.io"
	githubClientID = "Ov23liH6JRWDuToEcvw5"
)

// runRegister implements the register command using GitHub Device Flow
func runRegister(args []string) error {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Allow override via environment for testing
	baseURL := os.Getenv("BSUBIO_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	fmt.Println("Registering with bsub.io using GitHub authentication...")
	fmt.Println()

	// Step 1: Request device code
	deviceCode, userCode, verificationURI, expiresIn, interval, err := requestDeviceCode(baseURL, hostname)
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}

	// Step 2: Display code and prompt user
	fmt.Printf("! First copy your one-time code: %s\n", userCode)
	fmt.Printf("Press Enter to open %s in your browser...", verificationURI)

	// Wait for user to press Enter
	fmt.Scanln()

	// Open browser
	if err := openBrowser(verificationURI); err != nil {
		fmt.Printf("\nCould not open browser automatically. Please visit the URL above manually.\n")
	}

	fmt.Println("Waiting for authorization...")

	// Step 3: Poll for authorization
	apiKey, userInfo, err := pollForAuthorization(baseURL, deviceCode, userCode, interval, expiresIn)
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
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}

// requestDeviceCode initiates the device flow
func requestDeviceCode(baseURL, hostname string) (deviceCode, userCode, verificationURI string, expiresIn, interval int, err error) {
	url := fmt.Sprintf("%s/v1/auth/device/code", baseURL)

	reqBody := map[string]string{
		"hostname": hostname,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", "", 0, 0, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
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
func pollForAuthorization(baseURL, deviceCode, userCode string, interval, expiresIn int) (apiKey string, userInfo *UserInfo, err error) {
	url := fmt.Sprintf("%s/v1/auth/device/token", baseURL)
	pollInterval := time.Duration(interval) * time.Second
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)

	client := &http.Client{Timeout: 10 * time.Second}

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

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return "", nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", nil, err
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
