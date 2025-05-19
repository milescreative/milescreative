package routes

import (
	"net/http"

	"go-std/internal/utils"
)

type Router struct {
	mux *http.ServeMux
}

type Route struct {
	mux      *http.ServeMux
	path     string
	handlers map[string]http.HandlerFunc
}

func NewRouter(mux *http.ServeMux) *Router {
	return &Router{mux: mux}
}

func (r *Router) ROUTE(path string) *Route {
	route := &Route{
		mux:      r.mux,
		path:     path,
		handlers: make(map[string]http.HandlerFunc),
	}

	r.mux.HandleFunc(path, route.dispatch)
	return route
}

// dispatch is the single handler registered with ServeMux, it dispatches based on HTTP method
func (ri *Route) dispatch(w http.ResponseWriter, r *http.Request) {
	handler, ok := ri.handlers[r.Method]
	if !ok {
		utils.MethodNotAllowed(w, "Method not allowed")
		return
	}

	handler(w, r)
}

func (ri *Route) GET(handler http.HandlerFunc) *Route {
	ri.handlers[http.MethodGet] = _get(handler)
	return ri
}

func (ri *Route) POST(handler http.HandlerFunc) *Route {
	ri.handlers[http.MethodPost] = _post(handler)
	return ri
}

func (ri *Route) DELETE(handler http.HandlerFunc) *Route {
	ri.handlers[http.MethodDelete] = _delete(handler)
	return ri
}

func _get(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != "" {
			utils.MethodNotAllowed(w, "GET request should be used for this route.")
			//still continue on GET
		}
		next(w, r)
	})
}

func _post(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.MethodNotAllowed(w, "POST request required")
			return
		}
		next(w, r)
	})
}

func _delete(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.MethodNotAllowed(w, "DELETE request required")
			return
		}
		next(w, r)
	})
}

func ROOT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			utils.NotFound(w, "Not found")
			return
		}
		next(w, r)
	})
}
