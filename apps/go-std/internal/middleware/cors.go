package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// Create a single CORS handler instance
var corsHandler = cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3001"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"X-CSRF-Token"},
	AllowCredentials: true,
	Debug:            false,
})

// CORS middleware that uses the pre-configured handler
func CORS(next http.Handler) http.Handler {
	return corsHandler.Handler(next)
}
