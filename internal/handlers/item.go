package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tefsi/internal/domain"

	"github.com/go-chi/chi"
)

type ItemService interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItemByID(ctx context.Context, id int) (*domain.Item, error)
	GetItems(ctx context.Context, filter *domain.Filter) (*[]domain.Item, error)
	DeleteItem(ctx context.Context, id int) error
}

type ItemHandler struct {
	service ItemService
}

func NewItemHandler(service ItemService) *ItemHandler {
	return &ItemHandler{service}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	log.Println("received createitem request")
	var item domain.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Printf("bad json received, %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateItem(r.Context(), &item)
	if err != nil {
		log.Printf("error occured in createitem service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created item '%s' with id '%d'", item.Title, item.ID)

	w.WriteHeader(http.StatusCreated)
}

func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	log.Println("receieved getitembyid request")
	idStr := chi.URLParam(r, "id")

	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid item ID: %s", idStr)
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetItemByID(r.Context(), itemID)
	if err != nil {
		log.Printf("error occured in getitembyid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("responded with item '%s' with id %d", item.Title, itemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	log.Println("received getitems request")
	categoryIDString := chi.URLParam(r, "category")
	var categoryID int
	if categoryIDString != "" {
		var err error
		categoryID, err = strconv.Atoi(categoryIDString)

		if err != nil {
			log.Printf("got invalid category ID '%s'", categoryIDString)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	} else {
		categoryID = 0
	}

	filter := domain.Filter{
		CategoryID: categoryID,
	}

	itemList, err := h.service.GetItems(r.Context(), &filter)
	if err != nil {
		log.Printf("error occured in getitems service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("responded with list of items")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*itemList)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	log.Println("received deleteitem request")
	idStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid item ID '%s'", idStr)
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteItem(r.Context(), itemID)
	if err != nil {
		log.Printf("error occured in deleteitem service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("deleted item with id %d", itemID)

	w.WriteHeader(http.StatusOK)
}
