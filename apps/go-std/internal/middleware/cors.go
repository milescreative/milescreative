package middleware

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Enabling CORS")
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:8080"}, // Only allow your Next.js app
			AllowedMethods: []string{"GET", "POST"},           // Only allow specific methods
			AllowedHeaders: []string{"Content-Type"},          // Only allow specific headers
			Debug:          true,                              // Enable debugging for development
		})
		c.Handler(next).ServeHTTP(w, r)
	})
}
