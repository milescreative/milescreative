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
			utils.ErrorResponse(w, http.StatusForbidden, "Session required for CSRF protection", "SESSION_REQUIRED")
			return
		}

		// Get CSRF token from header
		token := r.Header.Get(csrfHeaderName)
		if token == "" {
			utils.ErrorResponse(w, http.StatusForbidden, "CSRF token missing from header", "CSRF_TOKEN_MISSING")
			return
		}
		//TODO add logger to app context and provide it here
		log.Println("CSRF token from header:", token)

		// Get stored HMAC from cookie
		csrfCookie, err := r.Cookie(csrfCookieName)
		if err != nil {
			utils.ErrorResponse(w, http.StatusForbidden, "CSRF token missing", "CSRF_TOKEN_MISSING")
			return
		}

		// Verify CSRF token
		if !utils.VerifyCSRFToken(token, sessionCookie.Value, csrfCookie.Value) {
			utils.ErrorResponse(w, http.StatusForbidden, "CSRF token mismatch", "CSRF_TOKEN_MISMATCH")
			return
		}

		next.ServeHTTP(w, r)
	})
}
