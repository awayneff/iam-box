// app/api/response.go

package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response structure for consistent API responses
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      int    `json:"code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithJSON sends a JSON response with status code
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
		}
	}
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, ErrorResponse{
		Error: message,
		Code:  status,
	})
}

// RespondWithInternalError sends a generic 500 error (hides details from client)
func RespondWithInternalError(w http.ResponseWriter, err error) {
	log.Printf("Internal error: %v", err)
	RespondWithError(w, http.StatusInternalServerError, "Internal server error")
}

// RespondWithValidationError sends a 400 with field-specific errors
func RespondWithValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "Validation failed",
		"fields": errors,
	})
}

// RespondWithNotFound sends a 404 response
func RespondWithNotFound(w http.ResponseWriter, resource string) {
	RespondWithError(w, http.StatusNotFound, resource+" not found")
}

// RespondWithCreated sends a 201 response with optional data
func RespondWithCreated(w http.ResponseWriter, data interface{}) {
	RespondWithJSON(w, http.StatusCreated, data)
}

// RespondWithNoContent sends a 204 response
func RespondWithNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RespondWithOK sends a 200 response with optional data
func RespondWithOK(w http.ResponseWriter, data interface{}) {
	RespondWithJSON(w, http.StatusOK, data)
}

// RespondWithOK sends a 405 response with optional data
func RespondWithBadMethod(w http.ResponseWriter, data interface{}) {
	RespondWithJSON(w, http.StatusMethodNotAllowed, data)
}
