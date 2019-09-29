package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// ItemHandler return a handler for the Rest API of a item
type ItemHandler struct {
	Service *services.ItemService
}

// NewRouter export a router configured with a supplier's routes
func (h *ItemHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", rootHandler(h.findAllItems))
	r.Method(http.MethodPost, "/", rootHandler(h.createItem))

	// Subroutes:
	r.Route("/{itemID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findItemByID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateItemID))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteItemByID))
	})

	return r
}

func (h *ItemHandler) findAllItems(w http.ResponseWriter, r *http.Request) error {
	items, err := h.Service.FindAllItems()
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) error {
	var payload dtos.ItemDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	result, err := h.Service.CreateItem(&payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	supplier, err := h.Service.FindItemByID(result)
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

func (h *ItemHandler) findItemByID(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	supplier, err := h.Service.FindItemByID(itemID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if supplier == nil {
		return NewNotFoundError(nil, "Item Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *ItemHandler) updateItemID(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	var payload dtos.ItemDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	supplier, err := h.Service.UpdateItemByID(itemID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if supplier == nil {
		return NewNotFoundError(nil, "Item Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(supplier); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *ItemHandler) deleteItemByID(w http.ResponseWriter, r *http.Request) error {
	itemID := chi.URLParam(r, "itemID")
	result, err := h.Service.DeleteItemID(itemID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "Item Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
