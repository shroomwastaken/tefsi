package services

import (
	"context"

	"tefsi/internal/domain"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserCartByID(ctx context.Context, id int) (*[]domain.Item, error)
	DeleteUser(ctx context.Context, id int) error
	CheckUserByDomain(ctx context.Context, user *domain.User) error
	UserExists(ctx context.Context, login string) error
}

// Реализация сервиса
type UserService struct {
	repo UserRepository
}

func NewDefaultUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserCartByID(ctx context.Context, id int) (*[]domain.Item, error) {
	return s.repo.GetUserCartByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) CheckUserByDomain(ctx context.Context, user *domain.User) error {
	return s.repo.CheckUserByDomain(ctx, user)
}

func (s *UserService) UserExists(ctx context.Context, login string) error {
	return s.repo.UserExists(ctx, login)
}
