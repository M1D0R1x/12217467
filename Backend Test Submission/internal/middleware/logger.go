package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Logger represents a simple logging interface
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct {
	output *os.File
}

// NewLogger creates a new SimpleLogger instance
func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		output: os.Stdout,
	}
}

// formatLog formats a log message with fields
func (l *SimpleLogger) formatLog(level, msg string, fields map[string]interface{}) string {
	timestamp := time.Now().Format(time.RFC3339)
	logMsg := fmt.Sprintf("[%s] %s - %s", level, timestamp, msg)

	if len(fields) > 0 {
		logMsg += " - Fields: "
		for k, v := range fields {
			logMsg += fmt.Sprintf("%s=%v ", k, v)
		}
	}

	return logMsg
}

// Info logs an informational message
func (l *SimpleLogger) Info(msg string, fields map[string]interface{}) {
	fmt.Fprintln(l.output, l.formatLog("INFO", msg, fields))
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, fields map[string]interface{}) {
	fmt.Fprintln(l.output, l.formatLog("ERROR", msg, fields))
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, fields map[string]interface{}) {
	fmt.Fprintln(l.output, l.formatLog("DEBUG", msg, fields))
}

// LoggingMiddleware creates a middleware that logs HTTP requests
func LoggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response wrapper to capture the status code
			rw := newResponseWriter(w)

			// Process the request
			next.ServeHTTP(rw, r)

			// Log the request details
			duration := time.Since(start)

			logger.Info("HTTP Request", map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     rw.statusCode,
				"duration":   duration.String(),
				"user_agent": r.UserAgent(),
				"remote_ip":  r.RemoteAddr,
			})
		})
	}
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
