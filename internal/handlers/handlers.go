package handlers

import "tefsi/internal/domain"

type Auth interface {
	GetUserFromJWT(header string) (*domain.User, error)
}

type AllHandlers struct {
	UserHandler     *UserHandler
	ItemHandler     *ItemHandler
	OrderHandler    *OrderHandler
	CategoryHandler *CategoryHandler
}
