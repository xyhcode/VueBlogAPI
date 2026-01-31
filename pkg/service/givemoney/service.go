/*
 * @Description: GiveMoney service interface
 * @Author: Qwenjie
 * @Date: 2026-01-24
 * @LastEditTime: 2026-01-24
 * @LastEditors: Qwenjie
 */
package givemoney

import (
	"context"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
)

// CreateGiveMoneyParams defines parameters for creating a give money record
type CreateGiveMoneyParams struct {
	Nickname string `json:"nickname" binding:"required"`
	Figure   int    `json:"figure" binding:"required"`
}

// UpdateGiveMoneyParams defines parameters for updating a give money record
type UpdateGiveMoneyParams struct {
	Nickname string `json:"nickname" binding:"required"`
	Figure   int    `json:"figure" binding:"required"`
}

// GiveMoneyService defines the business logic interface for give money operations
type GiveMoneyService interface {
	// GetAllRecords gets all give money records (no authentication required)
	GetAllRecords(ctx context.Context) ([]*model.GiveMoney, error)

	// GetRecordsByPage gets give money records with pagination (no authentication required)
	GetRecordsByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.GiveMoney], error)

	// CreateRecord creates a new give money record (authentication required)
	CreateRecord(ctx context.Context, params CreateGiveMoneyParams) (*model.GiveMoney, error)

	// UpdateRecord updates an existing give money record (authentication required)
	UpdateRecord(ctx context.Context, id uint, params UpdateGiveMoneyParams) (*model.GiveMoney, error)

	// DeleteRecord deletes a give money record (authentication required)
	DeleteRecord(ctx context.Context, id uint) error
}

type giveMoneyService struct {
	giveMoneyRepo repository.GiveMoneyRepository
}

// NewGiveMoneyService creates a new give money service
func NewGiveMoneyService(giveMoneyRepo repository.GiveMoneyRepository) GiveMoneyService {
	return &giveMoneyService{
		giveMoneyRepo: giveMoneyRepo,
	}
}

// GetAllRecords implements GiveMoneyService
func (s *giveMoneyService) GetAllRecords(ctx context.Context) ([]*model.GiveMoney, error) {
	return s.giveMoneyRepo.FindAll(ctx)
}

// GetRecordsByPage implements GiveMoneyService
func (s *giveMoneyService) GetRecordsByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.GiveMoney], error) {
	return s.giveMoneyRepo.FindListByPage(ctx, page, pageSize)
}

// CreateRecord implements GiveMoneyService
func (s *giveMoneyService) CreateRecord(ctx context.Context, params CreateGiveMoneyParams) (*model.GiveMoney, error) {
	giveMoney := &model.GiveMoney{
		Nickname: params.Nickname,
		Figure:   params.Figure,
	}

	err := s.giveMoneyRepo.Create(ctx, giveMoney)
	if err != nil {
		return nil, err
	}

	return giveMoney, nil
}

// UpdateRecord implements GiveMoneyService
func (s *giveMoneyService) UpdateRecord(ctx context.Context, id uint, params UpdateGiveMoneyParams) (*model.GiveMoney, error) {
	// First check if record exists
	existing, err := s.giveMoneyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the fields
	existing.Nickname = params.Nickname
	existing.Figure = params.Figure

	err = s.giveMoneyRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteRecord implements GiveMoneyService
func (s *giveMoneyService) DeleteRecord(ctx context.Context, id uint) error {
	return s.giveMoneyRepo.Delete(ctx, id)
}
