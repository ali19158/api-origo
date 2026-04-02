package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/online-shop/internal/middleware"
	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	userID, _ := middleware.UserIDFromContext(r.Context())

// 	var req models.CreateOrderRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		writeError(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	order, err := h.svc.Create(r.Context(), userID, req)
// 	if err != nil {
// 		if errors.Is(err, service.ErrInsufficientStock) {
// 			writeError(w, http.StatusConflict, err.Error())
// 			return
// 		}
// 		writeError(w, http.StatusInternalServerError, "failed to create order")
// 		return
// 	}

// 	writeJSON(w, http.StatusCreated, order)
// }

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	orders, err := h.svc.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list orders")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	var body struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, body.Status); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update order status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
