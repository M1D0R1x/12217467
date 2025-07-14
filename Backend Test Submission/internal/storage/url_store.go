package storage

import (
	"errors"
	"sync"
	"time"

	"12217467/backend_test_submission/internal/models"
)

var (
	// ErrShortcodeNotFound is returned when a shortcode is not found
	ErrShortcodeNotFound = errors.New("shortcode not found")

	// ErrShortcodeExists is returned when a shortcode already exists
	ErrShortcodeExists = errors.New("shortcode already exists")

	// ErrShortcodeExpired is returned when a shortcode has expired
	ErrShortcodeExpired = errors.New("shortcode has expired")
)

// URLStore defines the interface for URL storage operations
type URLStore interface {
	// Create stores a new short URL
	Create(shortURL models.ShortURL) error

	// Get retrieves a short URL by its shortcode
	Get(shortcode string) (models.ShortURL, error)

	// Update updates an existing short URL
	Update(shortURL models.ShortURL) error

	// Delete removes a short URL
	Delete(shortcode string) error

	// RecordClick records a click event for a shortcode
	RecordClick(shortcode string, click models.Click) error

	// ShortcodeExists checks if a shortcode already exists
	ShortcodeExists(shortcode string) bool
}

// InMemoryURLStore implements URLStore with in-memory storage
type InMemoryURLStore struct {
	urls  map[string]models.ShortURL
	mutex sync.RWMutex
}

// NewURLStore creates a new InMemoryURLStore
func NewURLStore() *InMemoryURLStore {
	return &InMemoryURLStore{
		urls: make(map[string]models.ShortURL),
	}
}

// Create stores a new short URL
func (s *InMemoryURLStore) Create(shortURL models.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.urls[shortURL.ID]; exists {
		return ErrShortcodeExists
	}

	s.urls[shortURL.ID] = shortURL
	return nil
}

// Get retrieves a short URL by its shortcode
func (s *InMemoryURLStore) Get(shortcode string) (models.ShortURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	shortURL, exists := s.urls[shortcode]
	if !exists {
		return models.ShortURL{}, ErrShortcodeNotFound
	}

	// Check if the URL has expired
	if time.Now().After(shortURL.ExpiresAt) {
		return models.ShortURL{}, ErrShortcodeExpired
	}

	return shortURL, nil
}

// Update updates an existing short URL
func (s *InMemoryURLStore) Update(shortURL models.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.urls[shortURL.ID]; !exists {
		return ErrShortcodeNotFound
	}

	s.urls[shortURL.ID] = shortURL
	return nil
}

// Delete removes a short URL
func (s *InMemoryURLStore) Delete(shortcode string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.urls[shortcode]; !exists {
		return ErrShortcodeNotFound
	}

	delete(s.urls, shortcode)
	return nil
}

// RecordClick records a click event for a shortcode
func (s *InMemoryURLStore) RecordClick(shortcode string, click models.Click) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	shortURL, exists := s.urls[shortcode]
	if !exists {
		return ErrShortcodeNotFound
	}

	// Check if the URL has expired
	if time.Now().After(shortURL.ExpiresAt) {
		return ErrShortcodeExpired
	}

	// Update click count and add click data
	shortURL.Clicks++
	shortURL.ClickData = append(shortURL.ClickData, click)

	// Update the URL in the store
	s.urls[shortcode] = shortURL
	return nil
}

// ShortcodeExists checks if a shortcode already exists
func (s *InMemoryURLStore) ShortcodeExists(shortcode string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.urls[shortcode]
	return exists
}
