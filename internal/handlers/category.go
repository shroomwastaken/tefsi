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

type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategories(ctx context.Context) (*[]domain.Category, error)
	DeleteCategory(ctx context.Context, id int) error
}

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("received create category request")
	var category domain.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Printf("bad json received, %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateCategory(r.Context(), &category)
	if err != nil {
		log.Printf("error occured in createcategories service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created category '%s'", category.Title)

	w.WriteHeader(http.StatusCreated)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	log.Println("receieved getcategorybyid request")
	idStr := chi.URLParam(r, "id")

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid category ID: %s", idStr)
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategoryByID(r.Context(), categoryID)
	if err != nil {
		log.Printf("error occured in getcategorybyid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("responded with category '%s' with id %d", category.Title, category.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	log.Println("received getcategories request")
	categoryList, err := h.service.GetCategories(r.Context())
	if err != nil {
		log.Printf("error occured in getcategories service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("responded with list of categories")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*categoryList)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("received deletecategory request")
	idStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid category id '%s'", idStr)
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		log.Printf("error occured in deletecategory service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("deleted category with id %d", categoryID)

	w.WriteHeader(http.StatusOK)
}
