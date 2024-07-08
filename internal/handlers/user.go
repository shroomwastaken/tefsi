package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"

	"tefsi/internal/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserCartByID(ctx context.Context, id int) (*[]domain.ItemWithAmount, error)
	DeleteUser(ctx context.Context, id int) error
	CheckUserByDomain(ctx context.Context, user *domain.User) error
	UserExists(ctx context.Context, login string) error
	UserIsAdmin(ctx context.Context, login string) (bool, error)
}

// Обработчики HTTP запросов
type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	log.Println("received getuserbyid request")
	id := chi.URLParam(r, "id")

	userID := 0
	var err error
	if userID, err = strconv.Atoi(id); err != nil {
		log.Printf("got invalid user ID '%s'", id)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("getting user %d", userID)
	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("error occured in getuserbyid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("responded with user with id %d", user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("received createuser request")
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("bad json received: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("creating user with login = %s, password = %s, is_admin = %v", user.Login, user.Password, user.IsAdmin)
	err = h.service.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("error occured in createuser service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created user with id %d", user.ID)

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.CheckUserByDomain(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := jwt.MapClaims{
		"sub": user.Login,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte("some_secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Authorization", "Bearer "+t)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.LoginResponse{Token: t})
}

func (h *UserHandler) GetUserCartByID(w http.ResponseWriter, r *http.Request) {
	log.Println("received getusercartbyid request")
	idStr := chi.URLParam(r, "id")

	cartID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid user ID '%s'", idStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cartItems, err := h.service.GetUserCartByID(r.Context(), cartID)

	if err != nil {
		log.Printf("error occured in getusercartbyid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("responded with users (id: %d) cart", cartID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*cartItems)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("received deleteuser request")
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid user ID: %s", idStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Printf("error occured in deleteuser service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("deleted user with id %d", userID)

	w.WriteHeader(http.StatusOK)
}

func (m *UserHandler) UserRequired(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("entering UserRequired middleware")
		login, err := getUserLoginFromJWT(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = m.service.UserExists(r.Context(), login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("leaving UserRequired middleware")
		next(w, r)
	}
}

func (m *UserHandler) AdminRequired(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("entering AdminRequired middleware")
		login, err := getUserLoginFromJWT(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isAdimn, err := m.service.UserIsAdmin(r.Context(), login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isAdimn {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Printf("leaving AdminRequired middleware")
		next(w, r)
	}
}

func getUserLoginFromJWT(header string) (string, error) {
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
