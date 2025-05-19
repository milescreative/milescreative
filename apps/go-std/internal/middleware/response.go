package middleware

import (
	"go-std/internal/utils"
	"net/http"
)

type Route struct{}

func NewRoute() *Route {
	return &Route{}
}

func (r *Route) Get(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "" {
			utils.ErrorResponse(w, http.StatusMethodNotAllowed, "GET request should be used for this route.", "METHOD_NOT_ALLOWED")
			//still continue as GET
		}
		next(w, r)
	})
}

func (r *Route) Post(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			utils.BadRequest(w, "POST request required")
			return
		}
		next(w, r)
	})
}

func (r *Route) Delete(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			utils.BadRequest(w, "DELETE request required")
			return
		}
		next(w, r)
	})
}
