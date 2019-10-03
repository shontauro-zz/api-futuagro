package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// CropHandler return a handler for the Rest API of a crop
type CropHandler struct {
	Service *services.CropService
}

// NewRouter export a router configured with the crop routes
func (h *CropHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", rootHandler(h.findAllCrops))
	r.Method(http.MethodPost, "/", rootHandler(h.createCrop))

	// Subroutes:
	r.Route("/{cropID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findCropByID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateCropByID))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteCropByID))
	})

	return r
}

func (h *CropHandler) findAllCrops(w http.ResponseWriter, r *http.Request) error {
	results, err := h.Service.FindAllCrops()
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (h *CropHandler) createCrop(w http.ResponseWriter, r *http.Request) error {
	var payload dtos.CropDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	crop, err := h.Service.CreateCrop(&payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(crop); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CropHandler) findCropByID(w http.ResponseWriter, r *http.Request) error {
	cropID := chi.URLParam(r, "cropID")
	crop, err := h.Service.FindCropByID(cropID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if crop == nil {
		return NewNotFoundError(nil, "Crop Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(crop); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CropHandler) updateCropByID(w http.ResponseWriter, r *http.Request) error {
	cropID := chi.URLParam(r, "cropID")
	var payload dtos.CropDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	crop, err := h.Service.UpdateCropByID(cropID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if crop == nil {
		return NewNotFoundError(nil, "Crop Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(crop); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CropHandler) deleteCropByID(w http.ResponseWriter, r *http.Request) error {
	cropID := chi.URLParam(r, "cropID")
	result, err := h.Service.DeleteCropByID(cropID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "Crop Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
