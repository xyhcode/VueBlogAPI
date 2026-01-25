/*
 * @Description: GiveMoney repository interface
 * @Author: Qwenjie
 * @Date: 2026-01-24
 * @LastEditTime: 2026-01-24
 * @LastEditors: Qwenjie
 */
package repository

import (
	"context"

	"github.com/anzhiyu-c/anheyu-app/pkg/domain/model"
)

// GiveMoneyRepository defines the contract for give money data operations.
type GiveMoneyRepository interface {
	// Embed base interface to get FindByID, Create, Update, Delete methods
	BaseRepository[model.GiveMoney]
	
	// FindAll gets all give money records without pagination
	FindAll(ctx context.Context) ([]*model.GiveMoney, error)
	
	// FindListByPage gets give money records with pagination
	FindListByPage(ctx context.Context, page, pageSize int) (*PageResult[model.GiveMoney], error)
}