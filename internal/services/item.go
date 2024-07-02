package services

import (
	"context"

	"tefsi/internal/domain"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItemByID(ctx context.Context, id int) (*domain.Item, error)
	GetItems(ctx context.Context, filter domain.Filter) ([]domain.Item, error)
}
