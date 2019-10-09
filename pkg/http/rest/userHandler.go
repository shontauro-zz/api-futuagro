package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
)

// UserHandler return a handler for the Rest API of an user
type UserHandler struct {
	Service *services.UserService
}

// NewRouter export a router configured with user routes
func (h *UserHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", rootHandler(h.findAllUsers))
	r.Method(http.MethodPost, "/", rootHandler(h.signup))

	// Subroutes:
	r.Route("/{userID}", func(r chi.Router) {
		r.Method(http.MethodGet, "/", rootHandler(h.findUserByID))
		r.Method(http.MethodPut, "/", rootHandler(h.updateUserByID))
		r.Method(http.MethodDelete, "/", rootHandler(h.deleteUserByID))
	})

	return r
}

func (h *UserHandler) findAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.Service.FindAllUsers()
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func (h *UserHandler) signup(w http.ResponseWriter, r *http.Request) error {
	var payload dtos.UserDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	user, err := h.Service.Signup(&payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *UserHandler) findUserByID(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	user, err := h.Service.PopulateUserByID(userID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if user == nil {
		return NewNotFoundError(nil, "User Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return NewAPIError(err, 500, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *UserHandler) updateUserByID(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	var payload dtos.UserDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	user, err := h.Service.UpdateUserByID(userID, &payload)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if user == nil {
		return NewNotFoundError(nil, "User Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func (h *UserHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	result, err := h.Service.DeleteUser(userID)
	if err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	if result == false {
		return NewNotFoundError(nil, "User Not Found")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return nil
}
