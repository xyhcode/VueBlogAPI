/*
 * @Description: Essay repository implementation
 * @Author: Qwenjie
 * @Date: 2026-01-27
 * @LastEditTime: 2026-01-27
 * @LastEditors: Qwenjie
 */
package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/anzhiyu-c/anheyu-app/ent"
	"github.com/anzhiyu-c/anheyu-app/ent/essay"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
	"github.com/anzhiyu-c/anheyu-app/pkg/domain/repository"
)

// toModel converts ent.Essay to model.Essay
func toModelEssay(e *ent.Essay) *model.Essay {
	if e == nil {
		return nil
	}

	modelEssay := &model.Essay{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Content:   e.Content,
		Date:      e.Date,
		Images:    e.Images, // Now stored as string
		Link:      e.Link,
		DeletedAt: e.DeletedAt,
	}

	// Convert images string to JSON
	_ = modelEssay.ConvertImagesStringToJSON()

	return modelEssay
}

type essayRepository struct {
	client *ent.Client
}

// NewEssayRepository creates a new essay repository
func NewEssayRepository(client *ent.Client) repository.EssayRepository {
	return &essayRepository{
		client: client,
	}
}

// FindByID finds an essay record by ID
func (r *essayRepository) FindByID(ctx context.Context, id uint) (*model.Essay, error) {
	essay, err := r.client.Essay.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return toModelEssay(essay), nil
}

// Create creates a new essay record
func (r *essayRepository) Create(ctx context.Context, entity *model.Essay) error {
	// Convert images JSON to string for storage
	_ = entity.ConvertImagesJSONToString()

	essay, err := r.client.Essay.Create().
		SetContent(entity.Content).
		SetDate(entity.Date).
		SetImages(entity.Images). // Set as string
		SetLink(entity.Link).
		Save(ctx)

	if err != nil {
		return err
	}

	// Update the entity with the created data
	entity.ID = essay.ID
	entity.CreatedAt = essay.CreatedAt
	entity.UpdatedAt = essay.UpdatedAt
	entity.DeletedAt = essay.DeletedAt

	// Convert images string back to JSON for return
	_ = entity.ConvertImagesStringToJSON()

	return nil
}

// Update updates an existing essay record
func (r *essayRepository) Update(ctx context.Context, entity *model.Essay) error {
	// Convert images JSON to string for storage
	_ = entity.ConvertImagesJSONToString()

	_, err := r.client.Essay.UpdateOneID(entity.ID).
		SetContent(entity.Content).
		SetDate(entity.Date).
		SetImages(entity.Images). // Set as string
		SetLink(entity.Link).
		Save(ctx)

	if err != nil {
		return err
	}

	// Convert images string back to JSON for return
	_ = entity.ConvertImagesStringToJSON()

	return err
}

// Delete deletes an essay record by ID
func (r *essayRepository) Delete(ctx context.Context, id uint) error {
	return r.client.Essay.DeleteOneID(id).Exec(ctx)
}

// FindListByPage gets essay records with pagination
func (r *essayRepository) FindListByPage(ctx context.Context, page, pageSize int) (*repository.PageResult[model.Essay], error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	total, err := r.client.Essay.Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	// Get paginated results - ordered by date descending (most recent first)
	essays, err := r.client.Essay.Query().
		Offset(offset).
		Limit(pageSize).
		Order(essay.ByDate(sql.OrderDesc())).
		All(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to model
	items := make([]*model.Essay, len(essays))
	for i, e := range essays {
		items[i] = toModelEssay(e)
	}

	return &repository.PageResult[model.Essay]{
		Items: items,
		Total: int64(total),
	}, nil
}

// FindAll gets all essay records ordered by date descending (most recent first)
func (r *essayRepository) FindAll(ctx context.Context) ([]*model.Essay, error) {
	essays, err := r.client.Essay.Query().
		Order(essay.ByDate(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Essay, len(essays))
	for i, e := range essays {
		result[i] = toModelEssay(e)
	}

	return result, nil
}
