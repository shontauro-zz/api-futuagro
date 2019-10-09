package rest

import (
	"encoding/json"
	"net/http"

	"futuagro.com/pkg/domain/dtos"
	"futuagro.com/pkg/domain/services"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler return a handler API for authorizing users
type AuthHandler struct {
	Service *services.AuthService
}

// NewRouter export a router configured with user routes
func (h *AuthHandler) NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/login", rootHandler(h.login))

	return r
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
	var payload dtos.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return NewAPIError(nil, http.StatusBadRequest, http.StatusBadRequest, "Bad request : invalid JSON.")
	}

	user, err := h.Service.Login(&payload)

	if err != nil {
		if errors.Cause(err) == bcrypt.ErrMismatchedHashAndPassword || errors.Cause(err) == bcrypt.ErrHashTooShort {
			return NewUnauthorizedError(err, "Authentication failed. Wrong user or password.")
		}
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if user == nil {
		return NewAPIError(nil, http.StatusUnauthorized, http.StatusUnauthorized, "Authentication failed. Wrong user or password.")
	}

	user.HashedPassword = ""
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return NewAPIError(err, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
