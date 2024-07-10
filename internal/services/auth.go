package services

import (
	"context"

	"tefsi/internal/domain"
)

type AuthRepository interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewDefaultAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	return s.repo.GetUserByLogin(ctx, login)
}
