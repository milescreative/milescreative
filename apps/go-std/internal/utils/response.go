// internal/utils/response.go
package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an API error response
type Error struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// JSON writes a JSON response to the ResponseWriter
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success writes a success response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Error writes an error response and optionally logs unexpected errors
func ErrorResponse(w http.ResponseWriter, status int, message string, code string) {
	JSONResponse(w, status, Response{
		Success: false,
		Error: &Error{
			Message: message,
			Code:    code,
		},
	})
}

// Common error responses for expected errors
func BadRequest(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusBadRequest, message, "BAD_REQUEST")
}

func Unauthorized(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnauthorized, message, "UNAUTHORIZED")
}

func Forbidden(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusForbidden, message, "FORBIDDEN")
}

func NotFound(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, "NOT_FOUND")
}

// Unexpected error responses
func InternalServerError(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusInternalServerError, message, "INTERNAL_SERVER_ERROR")
}

func DatabaseError(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusInternalServerError, "Database operation failed", "DATABASE_ERROR")
}

func UnexpectedError(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusInternalServerError, message, "UNEXPECTED_ERROR")
}

func Redirect(w http.ResponseWriter, r *http.Request, targetPath string) {
	// Get the scheme and host from the request
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Create base URL
	baseURL := &url.URL{
		Scheme: scheme,
		Host:   r.Host,
	}

	// Join the base URL with the target path
	targetURL := baseURL.ResolveReference(&url.URL{Path: path.Join("/", targetPath)})

	http.Redirect(w, r, targetURL.String(), http.StatusTemporaryRedirect)
}
