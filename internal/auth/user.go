package auth

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"tefsi/internal/domain"
)

type AuthService interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

type Auth struct {
	service AuthService
}

func NewAuth(service AuthService) *Auth {
	return &Auth{service: service}
}

func (a *Auth) GetUserFromJWT(header string) (*domain.User, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("invalid header: %s", header)
		}
	}()

	if header == "" {
		return nil, fmt.Errorf("no header provided")
	}
	t := strings.Split(header, " ")[1]
	log.Printf("got token %s", t)

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("some_secret"), nil
	})

	log.Printf("token parsed")
	if err != nil {
		return nil, err
	}
	log.Printf("loading payload")
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: %v", token)
	}
	log.Printf("payload loaded")
	login := payload["sub"].(string)
	user, err := a.service.GetUserByLogin(context.Background(), login)
	if err != nil {
		return nil, err
	}
	return user, nil
}
