package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Constants for valid parameter values
const (
	// API endpoint
	LogAPIEndpoint = "http://20.244.56.144/evaIuation-service/Iogs"

	// Valid stacks
	StackBackend  = "backend"
	StackFrontend = "frontend"

	// Valid levels
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

// Valid backend packages
var validBackendPackages = map[string]bool{
	"cache":      true,
	"controller": true,
	"cron_job":   true,
	"db":         true,
	"domain":     true,
	"handler":    true,
	"repository": true,
	"route":      true,
	"service":    true,
}

// Valid frontend packages
var validFrontendPackages = map[string]bool{
	"component": true,
	"hook":      true,
	"page":      true,
	"store":     true,
	"util":      true,
}

// LogRequest represents the structure of the log request
type LogRequest struct {
	Stack   string `json:"stack"`
	Level   string `json:"level"`
	Package string `json:"package"`
	Message string `json:"message"`
}

// LogResponse represents the structure of the log response
type LogResponse struct {
	LogID   string `json:"logID"`
	Message string `json:"message"`
}

// Log sends a log message to the test server
// stack: "backend" or "frontend"
// level: "debug", "info", "warn", "error", or "fatal"
// pkg: package name (depends on stack)
// message: log message
func Log(stack, level, pkg, message string) error {
	// Validate stack
	if stack != StackBackend && stack != StackFrontend {
		return fmt.Errorf("invalid stack: %s, must be 'backend' or 'frontend'", stack)
	}

	// Validate level
	if level != LevelDebug && level != LevelInfo && level != LevelWarn && level != LevelError && level != LevelFatal {
		return fmt.Errorf("invalid level: %s, must be 'debug', 'info', 'warn', 'error', or 'fatal'", level)
	}

	// Validate package based on stack
	if stack == StackBackend {
		if !validBackendPackages[pkg] {
			return fmt.Errorf("invalid backend package: %s", pkg)
		}
	} else if stack == StackFrontend {
		if !validFrontendPackages[pkg] {
			return fmt.Errorf("invalid frontend package: %s", pkg)
		}
	}

	// Create log request
	logReq := LogRequest{
		Stack:   stack,
		Level:   level,
		Package: pkg,
		Message: message,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(logReq)
	if err != nil {
		return fmt.Errorf("error marshaling log request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", LogAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending log request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("log API returned non-OK status: %d", resp.StatusCode)
	}

	// Parse response
	var logResp LogResponse
	if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
		return fmt.Errorf("error decoding log response: %w", err)
	}

	// Check if log was created successfully
	if logResp.Message != "log created successfully" {
		return errors.New("log was not created successfully")
	}

	return nil
}

// LogWithRetry attempts to send a log message with retries
func LogWithRetry(stack, level, pkg, message string, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := Log(stack, level, pkg, message)
		if err == nil {
			return nil
		}
		lastErr = err
		// Wait before retrying (exponential backoff)
		time.Sleep(time.Duration(1<<uint(i)) * 100 * time.Millisecond)
	}
	return fmt.Errorf("failed to send log after %d retries: %w", maxRetries, lastErr)
}
