package services

import (
	"context"

	"tefsi/internal/domain"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItemByID(ctx context.Context, id int) (*domain.Item, error)
	GetItems(ctx context.Context, filter *domain.Filter) (*[]domain.Item, error)
}

type ItemService struct {
	repo ItemRepository
}

func NewDefaultItemService(repo ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetItemByID(ctx context.Context, id int) (*domain.Item, error) {
	return s.repo.GetItemByID(ctx, id)
}

func (s *ItemService) CreateItem(ctx context.Context, item *domain.Item) error {
	return s.repo.CreateItem(ctx, item)
}

func (s *ItemService) GetItems(ctx context.Context, filter *domain.Filter) (*[]domain.Item, error) {
	return s.repo.GetItems(ctx, filter)
}
