// internal/infrastructure/http/response.go
package http

import (
	"encoding/json"
	"net/http"
)

type ResponseBaseDTO struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Errors     []string    `json:"errors,omitempty"`
}

func NewSuccessResponse(data interface{}, message string, statusCode int) *ResponseBaseDTO {
	return &ResponseBaseDTO{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewErrorResponse(message string, statusCode int, errors ...string) *ResponseBaseDTO {
	return &ResponseBaseDTO{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

func RespondSuccess(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	response := NewSuccessResponse(data, message, statusCode)
	respondJSON(w, statusCode, response)
}

func RespondError(w http.ResponseWriter, message string, statusCode int, errors ...string) {
	response := NewErrorResponse(message, statusCode, errors...)
	respondJSON(w, statusCode, response)
}

func respondJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
