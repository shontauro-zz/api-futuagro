package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// VariantHandler return a handler for the Rest API of a item
type VariantHandler struct {
	Service *services.VariantService
}

// NewRouter export a router configured with a supplier's routes
func (h *VariantHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", rootHandler(h.createVariant))
	r.Method(http.MethodGet, "/", rootHandler(h.findVariantsByItemID))

	// Subroutes:
	r.Route("/{variantID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findOneVariantByItemID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateVariant))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteVariant))
	})

	return r
}

func (h *VariantHandler) findVariantsByItemID(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	variants, err := h.Service.FindVariantsByItemID(itemID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(variants); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (h *VariantHandler) createVariant(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	var payload dtos.VariantDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	result, err := h.Service.CreateVariant(itemID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	supplier, err := h.Service.FindVariantByID(result)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *VariantHandler) findOneVariantByItemID(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	variantID := chi.URLParam(r, "variantID")
	variant, err := h.Service.FindOneVariantByItemID(itemID, variantID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if variant == nil {
		return NewNotFoundError(nil, "Variant Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(variant); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *VariantHandler) updateVariant(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	variantID := chi.URLParam(r, "variantID")
	var payload dtos.VariantDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	supplier, err := h.Service.UpdateVariant(itemID, variantID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if supplier == nil {
		return NewNotFoundError(nil, "Variant Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *VariantHandler) deleteVariant(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	variantID := chi.URLParam(r, "variantID")
	result, err := h.Service.DeleteVariant(itemID, variantID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "Variant Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
