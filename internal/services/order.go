package services

import (
	"context"

	"tefsi/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id int) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	GetOrders(ctx context.Context) (*[]domain.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}

type OrderService struct {
	repo OrderRepository
}

func NewDefaultOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (*domain.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrders(ctx context.Context) (*[]domain.Order, error) {
	return s.repo.GetOrders(ctx)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *domain.Order) error {
	return nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	return s.repo.DeleteOrder(ctx, id)
}
