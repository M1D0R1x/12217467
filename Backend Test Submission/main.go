package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"12217467/backend_test_submission/internal/api"
	"12217467/backend_test_submission/internal/middleware"
	"12217467/backend_test_submission/internal/storage"
)

func main() {
	// Initialize logger
	logger := middleware.NewLogger()

	// Initialize storage
	urlStore := storage.NewURLStore()

	// Initialize API handlers
	handler := api.NewHandler(urlStore, logger)

	// Create router and register routes
	mux := http.NewServeMux()

	// Register API endpoints
	mux.HandleFunc("/shorturls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateShortURL(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/shorturls/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path != "/shorturls/" {
			// The handler will extract the shortcode from the path
			handler.GetURLStats(w, r)
		} else {
			http.Error(w, "Method not allowed or invalid path", http.StatusMethodNotAllowed)
		}
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Redirection endpoint - catch-all handler for shortcodes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path != "/" {
			// This handles all paths except the root path
			handler.RedirectURL(w, r)
		} else if r.URL.Path == "/" {
			// Serve the index.html file for the root path
			http.ServeFile(w, r, "static/index.html")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Apply logging middleware
	wrappedMux := middleware.LoggingMiddleware(logger)(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      wrappedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
