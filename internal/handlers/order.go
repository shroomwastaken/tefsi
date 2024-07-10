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

type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id int) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	GetOrders(ctx context.Context) (*[]domain.Order, error)
	DeleteOrder(ctx context.Context, id int) error
	GetOrdersByUserID(ctx context.Context, id int) (*[]domain.Order, error)
}

type OrderHandler struct {
	service OrderService
	auth    Auth
}

func NewOrderHandler(service OrderService, auth Auth) *OrderHandler {
	return &OrderHandler{service, auth}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("received createorder request")
	var order domain.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	log.Printf("recieved order: %v", order)
	if err != nil {
		log.Printf("bad json received: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateOrder(r.Context(), &order)
	if err != nil {
		log.Printf("error occured in createorder service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created order with id %d", order.ID)

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	log.Println("received getorderbyid request")
	idStr := chi.URLParam(r, "id")

	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid order ID '%s'", idStr)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), orderID)
	if err != nil {
		log.Printf("error occured in getorderbyid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("responded with order with id %d", order.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("received getorders request")
	orderList, err := h.service.GetOrders(r.Context())
	if err != nil {
		log.Printf("error occured in getorders service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("responded with orders list")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

func (h *OrderHandler) UpdateOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("received updateorders request")
	var order domain.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Printf("bad json received: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrder(r.Context(), &order)
	if err != nil {
		log.Println("error occured in updateorder service")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("updated order with id %d", order.ID)

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("received deleteorder request")
	idStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid order ID '%s'", idStr)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteOrder(r.Context(), orderID)
	if err != nil {
		log.Printf("error occured in delete order service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("deleted order with id %d", orderID)

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	log.Println("received getordersbyuserid request")
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("got invalid user ID '%s'", idStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	orderList, err := h.service.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("error occured in getordersbyuserid service: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("responded with orders list")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}
