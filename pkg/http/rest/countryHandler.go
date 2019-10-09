package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// CountryHandler return a handler for the Rest API of a city
type CountryHandler struct {
	Service *services.CountryService
}

// NewRouter export a router configured with a country's routes
func (h *CountryHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", rootHandler(h.findAllCountries))
	r.Method(http.MethodPost, "/", rootHandler(h.createCountry))

	// Subroutes:
	r.Route("/{countryID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findCountryByID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateCountryByID))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteCountryByID))
	})

	r.Route("/{countryID}/country-states", func(r chi.Router) {
		r.Method(http.MethodPost, "/", rootHandler(h.createState))
		r.Method(http.MethodPut, "/{stateID}", rootHandler(h.updateState))
		r.Method(http.MethodDelete, "/{stateID}", rootHandler(h.deleteState))
	})

	return r
}

func (h *CountryHandler) findAllCountries(w http.ResponseWriter, r *http.Request) error {
	results, err := h.Service.FindAllCountries()
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

func (h *CountryHandler) createCountry(w http.ResponseWriter, r *http.Request) error {
	var payload dtos.CountryDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	result, err := h.Service.CreateCountry(&payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	country, err := h.Service.FindCountryByID(result)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CountryHandler) findCountryByID(w http.ResponseWriter, r *http.Request) error {
	ID := chi.URLParam(r, "countryID")
	country, err := h.Service.FindCountryByID(ID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if country == nil {
		return NewNotFoundError(nil, "Country Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CountryHandler) updateCountryByID(w http.ResponseWriter, r *http.Request) error {
	ID := chi.URLParam(r, "countryID")
	var payload dtos.CountryDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	country, err := h.Service.UpdateCountryByID(ID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if country == nil {
		return NewNotFoundError(nil, "Country Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CountryHandler) deleteCountryByID(w http.ResponseWriter, r *http.Request) error {
	ID := chi.URLParam(r, "countryID")
	result, err := h.Service.DeleteCountryByID(ID)
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

func (h *CountryHandler) createState(w http.ResponseWriter, r *http.Request) error {
	countryID := chi.URLParam(r, "countryID")
	var payload dtos.CountryStateDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	country, err := h.Service.AddState(countryID, payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if country == nil {
		return NewNotFoundError(nil, "Country Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CountryHandler) updateState(w http.ResponseWriter, r *http.Request) error {
	countryID := chi.URLParam(r, "countryID")
	stateID := chi.URLParam(r, "stateID")
	var payload dtos.CountryStateDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	country, err := h.Service.UpdateState(countryID, stateID, payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if country == nil {
		return NewNotFoundError(nil, "CountryState Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *CountryHandler) deleteState(w http.ResponseWriter, r *http.Request) error {
	countryID := chi.URLParam(r, "countryID")
	stateID := chi.URLParam(r, "stateID")
	country, err := h.Service.DeleteState(countryID, stateID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if country == nil {
		return NewNotFoundError(nil, "CountryState Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(country); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
