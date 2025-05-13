package middleware

import (
	"go-std/internal/utils"
	"net/http"
)

func (a *MiddlewareContext) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valid, err := utils.ValidateSession(a.Queries, w, r)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
			return
		}
		if !valid {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *MiddlewareContext) Protected(handler http.HandlerFunc) http.HandlerFunc {
	return a.AuthMiddleware(handler).ServeHTTP
}
