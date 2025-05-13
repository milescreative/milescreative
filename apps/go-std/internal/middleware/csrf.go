package middleware

import (
	"go-std/internal/utils"
	"log"
	"net/http"
)

const (
	csrfCookieName = utils.CsrfCookieName
	csrfHeaderName = utils.CsrfHeaderName
)

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF check for GET, HEAD, OPTIONS requests
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Get session ID from session cookie
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Session required for CSRF protection", http.StatusForbidden)
			return
		}

		// Get CSRF token from header
		token := r.Header.Get(csrfHeaderName)
		if token == "" {
			http.Error(w, "CSRF token missing from header", http.StatusForbidden)
			return
		}
		//TODO add logger to app context and provide it here
		log.Println("CSRF token from header:", token)

		// Get stored HMAC from cookie
		csrfCookie, err := r.Cookie(csrfCookieName)
		if err != nil {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		// Verify CSRF token
		if !utils.VerifyCSRFToken(token, sessionCookie.Value, csrfCookie.Value) {
			http.Error(w, "CSRF token mismatch", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
