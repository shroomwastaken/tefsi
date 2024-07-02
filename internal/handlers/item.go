package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"tefsi/internal/domain"

	"github.com/go-chi/chi"
)

type ItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItemByID(ctx context.Context, id int) (*domain.Item, error)
	GetItems(ctx context.Context, filter domain.Filter) ([]domain.Item, error)
}

type ItemHandler struct {
	service ItemService
}

func NewItemHandler(service ItemService) *ItemHandler {
	return &ItemHandler{service}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateItem(r.Context(), &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetItemByID(r.Context(), itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category")
	_ = categoryID

	// TODO: get category from ID
	var category domain.Category

	filter := domain.Filter{
		Category: category,
	}

	itemList, err := h.service.GetItems(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemList)
}
