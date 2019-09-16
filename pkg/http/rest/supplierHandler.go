package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/models"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// SupplierHandler return a handler for the Rest API of a supplier
type SupplierHandler struct {
	Service *services.SupplierService
}

// NewRouter export a router configured with a supplier's routes
func (h *SupplierHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", rootHandler(h.findAllSuppliers))
	r.Method(http.MethodPost, "/", rootHandler(h.createSupplier))

	// Subroutes:
	r.Route("/{id}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findSupplierByID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateSupplierByID))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteSupplierByID))
	})

	return r
}

func (h *SupplierHandler) findAllSuppliers(w http.ResponseWriter, r *http.Request) error {
	suppliers, err := h.Service.FindAllSuppliers()
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(suppliers); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (h *SupplierHandler) createSupplier(w http.ResponseWriter, r *http.Request) error {
	var payload models.Supplier
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	result, err := h.Service.CreateSupplier(&payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	supplier, err := h.Service.FindSupplierByID(result)
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

func (h *SupplierHandler) findSupplierByID(w http.ResponseWriter, r *http.Request) error {
	supplierID := chi.URLParam(r, "id")
	supplier, err := h.Service.FindSupplierByID(supplierID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if supplier == nil {
		return NewNotFoundError(nil, "Supplier Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *SupplierHandler) updateSupplierByID(w http.ResponseWriter, r *http.Request) error {
	supplierID := chi.URLParam(r, "id")
	var payload models.Supplier
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	supplier, err := h.Service.UpdateSupplierByID(supplierID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if supplier == nil {
		return NewNotFoundError(nil, "Supplier Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *SupplierHandler) deleteSupplierByID(w http.ResponseWriter, r *http.Request) error {
	supplierID := chi.URLParam(r, "id")
	result, err := h.Service.DeleteSupplier(supplierID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "Supplier Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
