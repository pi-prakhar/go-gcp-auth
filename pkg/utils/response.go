package utils

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	Write(w http.ResponseWriter)
}

type SuccessResponse[T any] struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       T      `json:"data,omitempty"`
}

func (r *SuccessResponse[T]) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	json.NewEncoder(w).Encode(r)
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
	Error      string `json:"error,omitempty"`
}

func (r *ErrorResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	json.NewEncoder(w).Encode(r)
}
