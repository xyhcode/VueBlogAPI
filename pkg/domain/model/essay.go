package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Image represents an image in an essay
type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt,omitempty"`
}

// Essay represents the essay entity
type Essay struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Content    string     `json:"content"`
	Date       time.Time  `json:"date"`
	Images     string     `json:"-"`                // Store as JSON string internally
	ImagesJSON []Image    `json:"images,omitempty"` // Exposed as JSON array in API
	Link       string     `json:"link,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// Helper function to convert images string to JSON
func (e *Essay) ConvertImagesStringToJSON() error {
	if e.Images != "" {
		var images []Image
		err := json.Unmarshal([]byte(e.Images), &images)
		if err != nil {
			return fmt.Errorf("failed to unmarshal images: %w", err)
		}
		e.ImagesJSON = images
	}
	return nil
}

// Helper function to convert images JSON to string
func (e *Essay) ConvertImagesJSONToString() error {
	if len(e.ImagesJSON) > 0 {
		bytes, err := json.Marshal(e.ImagesJSON)
		if err != nil {
			return fmt.Errorf("failed to marshal images: %w", err)
		}
		e.Images = string(bytes)
	} else {
		e.Images = ""
	}
	return nil
}

// Helper function to parse date string to time.Time
func ParseDateString(dateStr string) (time.Time, error) {
	// Try different common date formats
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // ISO 8601 without timezone
		"2006-01-02 15:04:05", // Common datetime format
		"2006-01-02",          // Date only
		"Jan 2, 2006",         // US format
		"January 2, 2006",     // Full month name
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// If none of the formats work, return an error
	return time.Time{}, fmt.Errorf("unable to parse date string: %s", dateStr)
}
