package services

import (
	"context"

	"tefsi/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order domain.Order) error
	GetOrderByID(ctx context.Context, id int) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order domain.Order) error
	GetOrders(ctx context.Context) ([]domain.Order, error)
}
