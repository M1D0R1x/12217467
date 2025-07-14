# Logging Middleware

A reusable logging middleware package that sends log messages to a remote logging service.

## Features

- Validates log parameters according to specified constraints
- Sends logs to a remote logging service via HTTP
- Supports retry mechanism for handling network issues
- Provides clear error messages for invalid parameters

## Installation

```bash
# If using Go modules
go get github.com/yourusername/logging-middleware
```

## Usage

```go
package main

import (
    "fmt"
    "https://github.com/M1D0R1x/12217467/logging-middleware/logger"
)

func main() {
    // Basic usage
    err := logger.Log("backend", "error", "handler", "received string, expected bool")
    if err != nil {
        fmt.Printf("Error sending log: %v\n", err)
    }

    // With retry mechanism
    err = logger.LogWithRetry("backend", "fatal", "db", "Critical database connection failure", 3)
    if err != nil {
        fmt.Printf("Error sending log after retries: %v\n", err)
    }
}
```

## API Reference

### `Log(stack, level, pkg, message string) error`

Sends a log message to the remote logging service.

Parameters:
- `stack` (string): The stack where the log originated. Must be one of:
  - `"backend"`
  - `"frontend"`
- `level` (string): The severity level of the log. Must be one of:
  - `"debug"`
  - `"info"`
  - `"warn"`
  - `"error"`
  - `"fatal"`
- `pkg` (string): The package where the log originated. 
  - For backend stack, must be one of: `"cache"`, `"controller"`, `"cron_job"`, `"db"`, `"domain"`, `"handler"`, `"repository"`, `"route"`, `"service"`
  - For frontend stack, must be one of: `"component"`, `"hook"`, `"page"`, `"store"`, `"util"`
- `message` (string): The log message.

Returns:
- `error`: An error if the log could not be sent or if any parameter is invalid.

### `LogWithRetry(stack, level, pkg, message string, maxRetries int) error`

Attempts to send a log message with retries in case of failure.

Parameters:
- Same as `Log()` plus:
- `maxRetries` (int): The maximum number of retry attempts.

Returns:
- `error`: An error if the log could not be sent after all retry attempts.

## Constants

The package provides constants for valid parameter values:

```go
// Stack constants
logger.StackBackend  // "backend"
logger.StackFrontend // "frontend"

// Level constants
logger.LevelDebug // "debug"
logger.LevelInfo  // "info"
logger.LevelWarn  // "warn"
logger.LevelError // "error"
logger.LevelFatal // "fatal"
```

## Error Handling

The `Log()` function returns an error in the following cases:
- Invalid stack parameter
- Invalid level parameter
- Invalid package parameter
- JSON marshaling error
- HTTP request creation error
- HTTP request sending error
- Non-OK response status
- Response parsing error
- Unsuccessful log creation

## Example Integration

Here's an example of how to integrate this logging middleware with an HTTP server:

```go
package main

import (
    "net/http"
    "https://github.com/M1D0R1x/12217467/logging-middleware/logger"
)

func main() {
    http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        // Log request received
        logger.Log("backend", "info", "handler", "Received request to /api/data")
        
        // Process request...
        
        // Log any errors
        if err := processRequest(); err != nil {
            logger.Log("backend", "error", "handler", "Error processing request: " + err.Error())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        
        // Log successful response
        logger.Log("backend", "info", "handler", "Successfully processed request")
        w.Write([]byte("Success"))
    })
    
    http.ListenAndServe(":8080", nil)
}
```