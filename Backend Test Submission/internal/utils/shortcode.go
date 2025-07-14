package utils

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
)

const (
	// DefaultShortcodeLength is the default length for generated shortcodes
	DefaultShortcodeLength = 6

	// MaxShortcodeLength is the maximum allowed length for custom shortcodes
	MaxShortcodeLength = 12
)

var (
	// ValidShortcodePattern defines the allowed characters in a shortcode
	ValidShortcodePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// GenerateShortcode creates a random shortcode of the specified length
func GenerateShortcode(length int) (string, error) {
	if length <= 0 {
		length = DefaultShortcodeLength
	}

	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64
	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim to desired length and remove any non-alphanumeric characters
	shortcode := encoded[:length]
	shortcode = strings.ReplaceAll(shortcode, "+", "a")
	shortcode = strings.ReplaceAll(shortcode, "/", "b")
	shortcode = strings.ReplaceAll(shortcode, "=", "c")

	return shortcode, nil
}

// ValidateShortcode checks if a shortcode is valid
func ValidateShortcode(shortcode string) bool {
	// Check if shortcode is empty
	if shortcode == "" {
		return false
	}

	// Check if shortcode is too long
	if len(shortcode) > MaxShortcodeLength {
		return false
	}

	// Check if shortcode contains only allowed characters
	return ValidShortcodePattern.MatchString(shortcode)
}
