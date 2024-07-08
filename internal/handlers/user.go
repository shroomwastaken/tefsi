package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"tefsi/internal/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserCartByID(ctx context.Context, id int) (*[]domain.ItemWithAmount, error)
	DeleteUser(ctx context.Context, id int) error
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

	err = h.service.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("error occured in createuser service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created user with id %d", user.ID)

	w.WriteHeader(http.StatusCreated)
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
