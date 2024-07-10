package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserLoginFromJWT(header string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("invalid header: %s", header)
		}
	}()
	if header == "" {
		return "", fmt.Errorf("no header provided")
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
		return "", err
	}
	log.Printf("loading payload")
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token: %v", token)
	}
	log.Printf("payload loaded")
	login := payload["sub"].(string)
	return login, nil
}
