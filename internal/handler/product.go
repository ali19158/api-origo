package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Create(r.Context(), &p); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create product")
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("failed to find product by id: %v", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.ProductFilter{
		Page:     intQuery(q.Get("page"), 1),
		PageSize: intQuery(q.Get("page_size"), 20),
	}

	if v := q.Get("category_id"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		filter.CategoryID = &id
	}
	if v := q.Get("min_price"); v != "" {
		f, _ := strconv.ParseFloat(v, 64)
		filter.MinPrice = &f
	}
	if v := q.Get("max_price"); v != "" {
		f, _ := strconv.ParseFloat(v, 64)
		filter.MaxPrice = &f
	}
	if v := q.Get("search"); v != "" {
		filter.Search = &v
	}
	if v := q.Get("is_featured"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.IsFeatured = &b
	}

	products, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to list products: %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"products":  products,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p.ID = id

	if err := h.svc.Update(r.Context(), &p); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func intQuery(s string, fallback int) int {
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}
