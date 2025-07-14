package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"12217467/backend_test_submission/internal/middleware"
	"12217467/backend_test_submission/internal/models"
	"12217467/backend_test_submission/internal/storage"
	"12217467/backend_test_submission/internal/utils"
)

const (
	// DefaultValidityMinutes is the default validity period in minutes
	DefaultValidityMinutes = 30
)

// Handler handles the API requests
type Handler struct {
	store  storage.URLStore
	logger middleware.Logger
}

// NewHandler creates a new Handler
func NewHandler(store storage.URLStore, logger middleware.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

// CreateShortURL handles the creation of a new short URL
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req models.CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate URL
	if req.URL == "" {
		h.respondWithError(w, http.StatusBadRequest, "URL is required", "")
		return
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid URL format", err.Error())
		return
	}

	// Set default validity if not provided
	validityMinutes := DefaultValidityMinutes
	if req.Validity != nil && *req.Validity > 0 {
		validityMinutes = *req.Validity
	}

	// Generate or validate shortcode
	shortcode := req.Shortcode
	if shortcode == "" {
		// Generate a random shortcode
		var err error
		shortcode, err = utils.GenerateShortcode(utils.DefaultShortcodeLength)
		if err != nil {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to generate shortcode", err.Error())
			return
		}
	} else {
		// Validate custom shortcode
		if !utils.ValidateShortcode(shortcode) {
			h.respondWithError(w, http.StatusBadRequest, "Invalid shortcode format", "Shortcode must be alphanumeric")
			return
		}

		// Check if shortcode already exists
		if h.store.ShortcodeExists(shortcode) {
			h.respondWithError(w, http.StatusConflict, "Shortcode already exists", "")
			return
		}
	}

	// Create short URL
	now := time.Now()
	shortURL := models.ShortURL{
		ID:          shortcode,
		OriginalURL: req.URL,
		CreatedAt:   now,
		ExpiresAt:   now.Add(time.Duration(validityMinutes) * time.Minute),
		Clicks:      0,
		ClickData:   []models.Click{},
	}

	// Store the short URL
	if err := h.store.Create(shortURL); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create short URL", err.Error())
		return
	}

	// Construct the short link
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	shortLink := fmt.Sprintf("%s://%s/%s", scheme, host, shortcode)

	// Prepare response
	resp := models.CreateShortURLResponse{
		ShortLink: shortLink,
		Expiry:    shortURL.ExpiresAt,
	}

	// Log success
	h.logger.Info("Created short URL", map[string]interface{}{
		"shortcode": shortcode,
		"url":       req.URL,
		"validity":  validityMinutes,
	})

	// Return response
	h.respondWithJSON(w, http.StatusCreated, resp)
}

// GetURLStats handles the retrieval of URL statistics
func (h *Handler) GetURLStats(w http.ResponseWriter, r *http.Request) {
	// Extract shortcode from path
	shortcode := strings.TrimPrefix(r.URL.Path, "/shorturls/")

	// Get URL from store
	shortURL, err := h.store.Get(shortcode)
	if err != nil {
		switch err {
		case storage.ErrShortcodeNotFound:
			h.respondWithError(w, http.StatusNotFound, "Shortcode not found", "")
		case storage.ErrShortcodeExpired:
			h.respondWithError(w, http.StatusGone, "Shortcode has expired", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve URL", err.Error())
		}
		return
	}

	// Prepare response
	resp := models.URLStatsResponse{
		Shortcode:   shortURL.ID,
		OriginalURL: shortURL.OriginalURL,
		CreatedAt:   shortURL.CreatedAt,
		ExpiresAt:   shortURL.ExpiresAt,
		Clicks:      shortURL.Clicks,
		ClickData:   shortURL.ClickData,
	}

	// Log success
	h.logger.Info("Retrieved URL stats", map[string]interface{}{
		"shortcode": shortcode,
		"clicks":    shortURL.Clicks,
	})

	// Return response
	h.respondWithJSON(w, http.StatusOK, resp)
}

// RedirectURL handles the redirection to the original URL
func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Extract shortcode from path
	shortcode := strings.TrimPrefix(r.URL.Path, "/")

	// Get URL from store
	shortURL, err := h.store.Get(shortcode)
	if err != nil {
		switch err {
		case storage.ErrShortcodeNotFound:
			h.respondWithError(w, http.StatusNotFound, "Shortcode not found", "")
		case storage.ErrShortcodeExpired:
			h.respondWithError(w, http.StatusGone, "Shortcode has expired", "")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve URL", err.Error())
		}
		return
	}

	// Record click
	click := models.Click{
		Timestamp: time.Now(),
		Referrer:  r.Referer(),
		Location:  getLocationFromIP(r.RemoteAddr),
		UserAgent: r.UserAgent(),
	}

	// Update click statistics asynchronously to not block the redirection
	go func() {
		if err := h.store.RecordClick(shortcode, click); err != nil {
			h.logger.Error("Failed to record click", map[string]interface{}{
				"shortcode": shortcode,
				"error":     err.Error(),
			})
		}
	}()

	// Log redirection
	h.logger.Info("Redirecting to original URL", map[string]interface{}{
		"shortcode": shortcode,
		"url":       shortURL.OriginalURL,
	})

	// Redirect to original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusFound)
}

// respondWithJSON sends a JSON response
func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// respondWithError sends an error response
func (h *Handler) respondWithError(w http.ResponseWriter, status int, message, details string) {
	errorResponse := models.ErrorResponse{
		Error:   message,
		Code:    status,
		Details: details,
	}

	h.logger.Error("API error", map[string]interface{}{
		"status":  status,
		"message": message,
		"details": details,
	})

	h.respondWithJSON(w, status, errorResponse)
}

// getLocationFromIP extracts a coarse-grained location from an IP address
// In a real application, this would use a geolocation service
func getLocationFromIP(ip string) string {
	// For simplicity, just return the IP address
	// In a production environment, this would use a geolocation service
	return fmt.Sprintf("Location from %s", ip)
}
