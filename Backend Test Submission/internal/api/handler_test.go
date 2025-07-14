package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"12217467/backend_test_submission/internal/models"
	"12217467/backend_test_submission/internal/storage"
)

// Logger interface for testing
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

// MockLogger is a simple mock implementation of the Logger interface for testing
type MockLogger struct{}

func (l *MockLogger) Info(msg string, fields map[string]interface{})  {}
func (l *MockLogger) Error(msg string, fields map[string]interface{}) {}
func (l *MockLogger) Debug(msg string, fields map[string]interface{}) {}

func TestCreateShortURL(t *testing.T) {
	// Setup
	store := storage.NewURLStore()
	logger := &MockLogger{}
	handler := NewHandler(store, logger)

	// Test case: Valid request with custom shortcode
	t.Run("Valid request with custom shortcode", func(t *testing.T) {
		// Create request
		reqBody := models.CreateShortURLRequest{
			URL:       "https://example.com",
			Validity:  intPtr(60),
			Shortcode: "testcode",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/shorturls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateShortURL(w, req)

		// Check response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var resp models.CreateShortURLResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify shortcode in response
		expectedShortLink := "http://" + req.Host + "/testcode"
		if resp.ShortLink != expectedShortLink {
			t.Errorf("Expected shortLink %s, got %s", expectedShortLink, resp.ShortLink)
		}
	})

	// Test case: Valid request without custom shortcode
	t.Run("Valid request without custom shortcode", func(t *testing.T) {
		// Create request
		reqBody := models.CreateShortURLRequest{
			URL:      "https://example.com",
			Validity: intPtr(60),
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/shorturls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateShortURL(w, req)

		// Check response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var resp models.CreateShortURLResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify that a shortLink was generated
		if resp.ShortLink == "" {
			t.Error("Expected non-empty shortLink")
		}
	})

	// Test case: Invalid URL
	t.Run("Invalid URL", func(t *testing.T) {
		// Create request with invalid URL
		reqBody := models.CreateShortURLRequest{
			URL: "not-a-valid-url",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/shorturls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateShortURL(w, req)

		// Check response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test case: Duplicate shortcode
	t.Run("Duplicate shortcode", func(t *testing.T) {
		// First request to create the shortcode
		reqBody := models.CreateShortURLRequest{
			URL:       "https://example.com",
			Shortcode: "duplicate",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/shorturls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateShortURL(w, req)

		// Second request with the same shortcode
		req = httptest.NewRequest("POST", "/shorturls", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w = httptest.NewRecorder()

		// Call handler
		handler.CreateShortURL(w, req)

		// Check response
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status code %d, got %d", http.StatusConflict, w.Code)
		}
	})
}

func TestGetURLStats(t *testing.T) {
	// Setup
	store := storage.NewURLStore()
	logger := &MockLogger{}
	handler := NewHandler(store, logger)

	// Create a test URL
	shortcode := "teststats"
	now := time.Now()
	shortURL := models.ShortURL{
		ID:          shortcode,
		OriginalURL: "https://example.com",
		CreatedAt:   now,
		ExpiresAt:   now.Add(30 * time.Minute),
		Clicks:      0,
		ClickData:   []models.Click{},
	}
	store.Create(shortURL)

	// Test case: Get stats for existing shortcode
	t.Run("Get stats for existing shortcode", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest("GET", "/shorturls/"+shortcode, nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.GetURLStats(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var resp models.URLStatsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify response data
		if resp.Shortcode != shortcode {
			t.Errorf("Expected shortcode %s, got %s", shortcode, resp.Shortcode)
		}
		if resp.OriginalURL != "https://example.com" {
			t.Errorf("Expected originalUrl %s, got %s", "https://example.com", resp.OriginalURL)
		}
	})

	// Test case: Get stats for non-existent shortcode
	t.Run("Get stats for non-existent shortcode", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest("GET", "/shorturls/nonexistent", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.GetURLStats(w, req)

		// Check response
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestRedirectURL(t *testing.T) {
	// Setup
	store := storage.NewURLStore()
	logger := &MockLogger{}
	handler := NewHandler(store, logger)

	// Create a test URL
	shortcode := "testredirect"
	now := time.Now()
	shortURL := models.ShortURL{
		ID:          shortcode,
		OriginalURL: "https://example.com",
		CreatedAt:   now,
		ExpiresAt:   now.Add(30 * time.Minute),
		Clicks:      0,
		ClickData:   []models.Click{},
	}
	store.Create(shortURL)

	// Test case: Redirect for existing shortcode
	t.Run("Redirect for existing shortcode", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest("GET", "/"+shortcode, nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.RedirectURL(w, req)

		// Check response
		if w.Code != http.StatusFound {
			t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
		}

		// Check redirect location
		location := w.Header().Get("Location")
		if location != "https://example.com" {
			t.Errorf("Expected redirect to %s, got %s", "https://example.com", location)
		}
	})

	// Test case: Redirect for non-existent shortcode
	t.Run("Redirect for non-existent shortcode", func(t *testing.T) {
		// Create request
		req := httptest.NewRequest("GET", "/nonexistent", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.RedirectURL(w, req)

		// Check response
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

// Helper function to create a pointer to an int
func intPtr(i int) *int {
	return &i
}
