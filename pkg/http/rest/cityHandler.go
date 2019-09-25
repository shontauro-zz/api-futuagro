package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// CityHandler return a handler for the Rest API of a city
type CityHandler struct {
	Service *services.CityService
}

// NewRouter export a router configured with a country's routes
func (h *CityHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/{stateID}/cities", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findAllCitiesByState))
		r.Method(http.MethodPost, "/", rootHandler(h.createCity))
		r.Method(http.MethodPut, "/{cityID}", rootHandler(h.updateCityByID))
		r.Method(http.MethodDelete, "/{cityID}", rootHandler(h.deleteCityByID))
	})

	return r
}

func (h *CityHandler) findAllCitiesByState(w http.ResponseWriter, r *http.Request) error {
	stateID := chi.URLParam(r, "stateID")
	results, err := h.Service.FindAllCitiesByCountryState(stateID)
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

func (h *CityHandler) createCity(w http.ResponseWriter, r *http.Request) error {
	stateID := chi.URLParam(r, "stateID")
	var payload dtos.CityDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	result, err := h.Service.CreateCity(stateID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	city, err := h.Service.FindCityByID(result)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(city); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CityHandler) findCityByID(w http.ResponseWriter, r *http.Request) error {
	ID := chi.URLParam(r, "id")
	city, err := h.Service.FindCityByID(ID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if city == nil {
		return NewNotFoundError(nil, "City Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(city); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CityHandler) updateCityByID(w http.ResponseWriter, r *http.Request) error {
	stateID := chi.URLParam(r, "stateID")
	cityID := chi.URLParam(r, "cityID")
	var payload dtos.CityDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	city, err := h.Service.UpdateCityByID(stateID, cityID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if city == nil {
		return NewNotFoundError(nil, "City Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(city); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CityHandler) deleteCityByID(w http.ResponseWriter, r *http.Request) error {
	stateID := chi.URLParam(r, "stateID")
	cityID := chi.URLParam(r, "cityID")
	result, err := h.Service.DeleteCityByID(stateID, cityID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "Country Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
