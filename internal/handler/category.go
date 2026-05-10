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

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Create(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	cat, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("category not found %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, cat)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to list categories: %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, cats)
}

// ListRootCategories handles GET /api/v2/categories — only root categories (parent_id IS NULL).
func (h *CategoryHandler) ListRootCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListRootCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to list root categories: %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, cats)
}

// ListSubcategories handles GET /api/v1/categories/{category_id}/subcategories — direct children.
func (h *CategoryHandler) ListSubcategories(w http.ResponseWriter, r *http.Request) {
	parentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	cats, err := h.svc.ListSubcategories(r.Context(), parentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to list subcategories: %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, cats)
}

// ListTree handles GET /api/v1/categories/tree — full category tree with nested children.
func (h *CategoryHandler) ListTree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.svc.ListTree(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to build category tree: %w", err).Error())
		return
	}

	writeJSON(w, http.StatusOK, tree)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c.ID = id

	if err := h.svc.Update(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update category")
		return
	}

	writeJSON(w, http.StatusOK, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
