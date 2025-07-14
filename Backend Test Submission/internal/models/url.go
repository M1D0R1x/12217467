package models

import (
	"time"
)

// ShortURL represents a shortened URL with its metadata
type ShortURL struct {
	ID          string    `json:"id"`          // Unique identifier (shortcode)
	OriginalURL string    `json:"originalUrl"` // Original long URL
	CreatedAt   time.Time `json:"createdAt"`   // Creation timestamp
	ExpiresAt   time.Time `json:"expiresAt"`   // Expiration timestamp
	Clicks      int       `json:"clicks"`      // Number of times the URL has been accessed
	ClickData   []Click   `json:"clickData"`   // Detailed click data
}

// Click represents a single click event on a shortened URL
type Click struct {
	Timestamp time.Time `json:"timestamp"` // When the click occurred
	Referrer  string    `json:"referrer"`  // Where the click came from
	Location  string    `json:"location"`  // Approximate geographical location
	UserAgent string    `json:"userAgent"` // User agent of the client
}

// CreateShortURLRequest represents the request body for creating a short URL
type CreateShortURLRequest struct {
	URL       string `json:"url"`       // Original URL to shorten
	Validity  *int   `json:"validity"`  // Optional validity period in minutes
	Shortcode string `json:"shortcode"` // Optional custom shortcode
}

// CreateShortURLResponse represents the response for a successful short URL creation
type CreateShortURLResponse struct {
	ShortLink string    `json:"shortLink"` // The complete shortened URL
	Expiry    time.Time `json:"expiry"`    // Expiration timestamp
}

// URLStatsResponse represents the response for URL statistics
type URLStatsResponse struct {
	Shortcode   string    `json:"shortcode"`   // The shortcode
	OriginalURL string    `json:"originalUrl"` // Original long URL
	CreatedAt   time.Time `json:"createdAt"`   // Creation timestamp
	ExpiresAt   time.Time `json:"expiresAt"`   // Expiration timestamp
	Clicks      int       `json:"clicks"`      // Total number of clicks
	ClickData   []Click   `json:"clickData"`   // Detailed click data
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`   // Error message
	Code    int    `json:"code"`    // HTTP status code
	Details string `json:"details"` // Additional error details (optional)
}