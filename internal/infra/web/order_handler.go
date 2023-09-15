package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/renan5g/go-clean-arch/internal/application/usecase"
)

type WebOrderHandler struct {
	CreateOrder usecase.CreateOrderUseCase
	ListOrder   usecase.ListOrderUseCase
}

func NewWebOrderHandler(
	r chi.Router,
	createOrder usecase.CreateOrderUseCase,
	listOrder usecase.ListOrderUseCase,
) *WebOrderHandler {
	handler := &WebOrderHandler{
		CreateOrder: createOrder,
		ListOrder:   listOrder,
	}

	r.Post("/orders", handler.Create)
	r.Get("/orders", handler.List)
	return handler
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrder.Execute(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrder.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
