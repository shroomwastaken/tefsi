package services

import (
	"context"

	"tefsi/internal/domain"
)

type StatusRepository interface {
	CreateStatus(ctx context.Context, status *domain.Status) error
	GetStatusByID(ctx context.Context, id int) (*domain.Status, error)
}
