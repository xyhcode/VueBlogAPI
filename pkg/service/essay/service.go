/*
 * @Description: Essay service interface
 * @Author: Qwenjie
 * @Date: 2026-01-27
 * @LastEditTime: 2026-01-27
 * @LastEditors: Qwenjie
 */
package essay

import (
	"context"
	"fmt"
	"time"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
)

// CreateEssayParams defines parameters for creating an essay record
type CreateEssayParams struct {
	Content string        `json:"content" binding:"required"`
	Date    string        `json:"date" binding:"required"`
	Images  []model.Image `json:"images,omitempty"` // JSON array
	Link    string        `json:"link,omitempty"`
}

// UpdateEssayParams defines parameters for updating an essay record
type UpdateEssayParams struct {
	Content string        `json:"content" binding:"required"`
	Date    string        `json:"date" binding:"required"`
	Images  []model.Image `json:"images,omitempty"` // JSON array
	Link    string        `json:"link,omitempty"`
}

// Service defines the business logic interface for essay operations
type Service interface {
	// GetAllEssays gets all essay records (no authentication required)
	GetAllEssays(ctx context.Context) ([]*model.Essay, error)

	// GetEssaysByPage gets essay records with pagination (no authentication required)
	GetEssaysByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.Essay], error)

	// CreateEssay creates a new essay record (authentication required)
	CreateEssay(ctx context.Context, params CreateEssayParams) (*model.Essay, error)

	// UpdateEssay updates an existing essay record (authentication required)
	UpdateEssay(ctx context.Context, id uint, params UpdateEssayParams) (*model.Essay, error)

	// DeleteEssay deletes an essay record (authentication required)
	DeleteEssay(ctx context.Context, id uint) error
}

type essayService struct {
	essayRepo repository.EssayRepository
}

// NewService creates a new essay service
func NewService(essayRepo repository.EssayRepository) Service {
	return &essayService{
		essayRepo: essayRepo,
	}
}

// GetAllEssays implements Service
func (s *essayService) GetAllEssays(ctx context.Context) ([]*model.Essay, error) {
	return s.essayRepo.FindAll(ctx)
}

// GetEssaysByPage implements Service
func (s *essayService) GetEssaysByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.Essay], error) {
	return s.essayRepo.FindListByPage(ctx, page, pageSize)
}

// CreateEssay implements Service
func (s *essayService) CreateEssay(ctx context.Context, params CreateEssayParams) (*model.Essay, error) {
	// Parse date string to time.Time
	date, err := parseDateString(params.Date)
	if err != nil {
		return nil, err
	}

	essay := &model.Essay{
		Content:    params.Content,
		Date:       date,
		ImagesJSON: params.Images, // Store as JSON array internally
		Link:       params.Link,
	}

	err = s.essayRepo.Create(ctx, essay)
	if err != nil {
		return nil, err
	}

	return essay, nil
}

// UpdateEssay implements Service
func (s *essayService) UpdateEssay(ctx context.Context, id uint, params UpdateEssayParams) (*model.Essay, error) {
	// First check if record exists
	existing, err := s.essayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Parse date string to time.Time
	date, err := parseDateString(params.Date)
	if err != nil {
		return nil, err
	}

	// Update the fields
	existing.Content = params.Content
	existing.Date = date
	existing.ImagesJSON = params.Images // Store as JSON array
	existing.Link = params.Link

	err = s.essayRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteEssay implements Service
func (s *essayService) DeleteEssay(ctx context.Context, id uint) error {
	return s.essayRepo.Delete(ctx, id)
}

// Helper function to parse date string to time.Time
func parseDateString(dateStr string) (time.Time, error) {
	// Try different common date formats
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // ISO 8601 without timezone
		"2006-01-02 15:04:05", // Common datetime format
		"2006-01-02",          // Date only
		"Jan 2, 2006",         // US format
		"January 2, 2006",     // Full month name,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// If none of the formats work, return an error
	return time.Time{}, fmt.Errorf("unable to parse date string: %s", dateStr)
}
