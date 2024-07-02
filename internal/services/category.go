package services

import (
	"context"

	"tefsi/internal/domain"
)

type CategoryRepository interface {
	CreateGategory(ctx context.Context, category domain.Category) error
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}
