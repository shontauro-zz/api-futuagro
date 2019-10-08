package rest

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// APIError implements ClientError interface.
type APIError struct {
	Cause   error  `json:"-"`
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

// ResponseBody returns JSON response body.
func (e *APIError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrap(err, "Error while parsing response body")
	}
	return body, nil
}

// ResponseHeaders returns http status code and headers.
func (e *APIError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// NewAPIError create an instance for an API error
func NewAPIError(err error, status int, code int, message string) error {
	return &APIError{
		Cause:   err,
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// NewNotFoundError create an error instance for an http error 404
func NewNotFoundError(err error, message string) error {
	return &APIError{
		Cause:   err,
		Status:  http.StatusNotFound,
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// NewUnauthorizedError create an error instance for an http error 401
func NewUnauthorizedError(err error, message string) error {
	return &APIError{
		Cause:   err,
		Status:  http.StatusUnauthorized,
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
