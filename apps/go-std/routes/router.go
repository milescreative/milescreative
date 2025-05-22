package routes

import (
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

type Router_ struct {
	mux *http.ServeMux
}

type Router struct {
	mux *httprouter.Router
}

func NewRouter_(mux *http.ServeMux) *Router_ {
	return &Router_{mux: mux}
}

func NewRouter() *Router {
	return &Router{mux: httprouter.New()}
}

var re = regexp.MustCompile(`\{(\w+)\}`)

// wraps a httprouter function to make it more like std
func w(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for _, param := range p {
			r.SetPathValue(param.Key, param.Value)
		}
		h.ServeHTTP(w, r)
	}
}

func (r *Router) GET_(path string, handler http.HandlerFunc) {
	// newpath := re.ReplaceAllString(path, `:$1`)
	r.mux.GET(path, w(handler))
}

func (r *Router) POST_(path string, handler http.HandlerFunc) {
	// newpath := re.ReplaceAllString(path, `:$1`)
	r.mux.POST(path, w(handler))
}

func (r *Router) PUT_(path string, handler http.HandlerFunc) {
	// newpath := re.ReplaceAllString(path, `:$1`)
	r.mux.PUT(path, w(handler))
}

func (r *Router) DELETE_(path string, handler http.HandlerFunc) {
	// newpath := re.ReplaceAllString(path, `:$1`)
	r.mux.DELETE(path, w(handler))
}

func (r *Router_) GET(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc("GET "+path, handler)
}

func (r *Router_) POST(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc("POST "+path, handler)
}

func (r *Router_) PUT(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc("PUT "+path, handler)
}

func (r *Router_) DELETE(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc("DELETE "+path, handler)
}
