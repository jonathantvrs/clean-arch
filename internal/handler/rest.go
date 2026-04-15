package handler

import (
	"encoding/json"
	"net/http"
	"order-service/internal/service"

	"github.com/gorilla/mux"
)

type RestHandler struct {
	service *service.OrderService
}

func NewRestHandler(s *service.OrderService) *RestHandler {
	return &RestHandler{service: s}
}

func (h *RestHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func SetupRoutes(r *mux.Router, s *service.OrderService) {
	h := NewRestHandler(s)
	r.HandleFunc("/orders", h.ListOrders).Methods("GET")
}
