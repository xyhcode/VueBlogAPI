/*
 * @Description: GiveMoney repository implementation
 * @Author: Qwenjie
 * @Date: 2026-01-24
 * @LastEditTime: 2026-01-24
 * @LastEditors: Qwenjie
 */
package ent

import (
	"context"

	"github.com/anzhiyu-c/anheyu-app/ent"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
)

// toModel converts ent.GiveMoney to model.GiveMoney
func toModel(gm *ent.GiveMoney) *model.GiveMoney {
	if gm == nil {
		return nil
	}
	
	return &model.GiveMoney{
		ID:        gm.ID,
		CreatedAt: gm.CreatedAt,
		UpdatedAt: gm.UpdatedAt,
		Nickname:  gm.Nickname,
		Figure:    gm.Figure,
		DeletedAt: gm.DeletedAt,
	}
}

type giveMoneyRepository struct {
	client *ent.Client
}

// NewGiveMoneyRepository creates a new give money repository
func NewGiveMoneyRepository(client *ent.Client) repository.GiveMoneyRepository {
	return &giveMoneyRepository{
		client: client,
	}
}

// FindByID finds a give money record by ID
func (r *giveMoneyRepository) FindByID(ctx context.Context, id uint) (*model.GiveMoney, error) {
	giveMoney, err := r.client.GiveMoney.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return toModel(giveMoney), nil
}

// Create creates a new give money record
func (r *giveMoneyRepository) Create(ctx context.Context, entity *model.GiveMoney) error {
	giveMoney, err := r.client.GiveMoney.Create().
		SetNickname(entity.Nickname).
		SetFigure(entity.Figure).
		Save(ctx)
	
	if err != nil {
		return err
	}
	
	// Update the entity with the created data
	entity.ID = giveMoney.ID
	entity.CreatedAt = giveMoney.CreatedAt
	entity.UpdatedAt = giveMoney.UpdatedAt
	entity.DeletedAt = giveMoney.DeletedAt
	
	return nil
}

// Update updates an existing give money record
func (r *giveMoneyRepository) Update(ctx context.Context, entity *model.GiveMoney) error {
	_, err := r.client.GiveMoney.UpdateOneID(entity.ID).
		SetNickname(entity.Nickname).
		SetFigure(entity.Figure).
		Save(ctx)
	
	return err
}

// Delete deletes a give money record by ID
func (r *giveMoneyRepository) Delete(ctx context.Context, id uint) error {
	return r.client.GiveMoney.DeleteOneID(id).Exec(ctx)
}

// FindAll gets all give money records without pagination
func (r *giveMoneyRepository) FindAll(ctx context.Context) ([]*model.GiveMoney, error) {
	giveMonies, err := r.client.GiveMoney.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	
	result := make([]*model.GiveMoney, len(giveMonies))
	for i, gm := range giveMonies {
		result[i] = toModel(gm)
	}
	
	return result, nil
}

// FindListByPage gets give money records with pagination
func (r *giveMoneyRepository) FindListByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.GiveMoney], error) {
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Get total count
	total, err := r.client.GiveMoney.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get paginated results
	giveMonies, err := r.client.GiveMoney.Query().
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	
	if err != nil {
		return nil, err
	}
	
	// Convert to model
	items := make([]*model.GiveMoney, len(giveMonies))
	for i, gm := range giveMonies {
		items[i] = toModel(gm)
	}
	
	return &repository.PageResult[model.GiveMoney]{
		Items: items,
		Total: int64(total),
	}, nil
}