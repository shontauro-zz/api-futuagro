package rest

import (
	"log"
	"net/http"
)

// Use rootHandler as wrapper around handler functions
type rootHandler func(http.ResponseWriter, *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	// Error handling
	log.Println(err)

	clientError, ok := err.(ClientError)
	if !ok {
		// If the error is not ClientError, assume that it is ServerError.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := clientError.ResponseBody()
	if err != nil {
		log.Printf("An error ocurred: %+v", err)
		w.WriteHeader(500)
		return
	}

	status, headers := clientError.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

// ClientError is an error whose details to be shared with client
type ClientError interface {
	Error() string
	// ResponseBody returns response body.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}
