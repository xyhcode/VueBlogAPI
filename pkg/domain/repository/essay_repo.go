/*
 * @Description: Essay repository interface
 * @Author: Qwenjie
 * @Date: 2026-01-27
 * @LastEditTime: 2026-01-27
 * @LastEditors: Qwenjie
 */
package repository

import (
	"context"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
)

// EssayRepository defines the contract for essay data operations.
type EssayRepository interface {
	// Embed base interface to get FindByID, Create, Update, Delete methods
	BaseRepository[model.Essay]

	// FindAll gets all essay records without pagination
	FindAll(ctx context.Context) ([]*model.Essay, error)

	// FindListByPage gets essay records with pagination
	FindListByPage(ctx context.Context, page, pageSize int) (*PageResult[model.Essay], error)
}
