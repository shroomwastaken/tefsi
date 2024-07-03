package services

import (
	"context"

	"tefsi/internal/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategories(ctx context.Context) (*[]domain.Category, error)
}

type CategoryService struct {
	repo CategoryRepository
}

func NewDefaultCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) GetCategories(ctx context.Context) (*[]domain.Category, error) {
	return s.repo.GetCategories(ctx)
}
