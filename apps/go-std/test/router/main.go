package main

import (
	"go-std/internal/config"

	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {

	env, _ := config.Config()
	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	port, _ := env.GetInt("APP_PORT")
	if port == 0 {
		port = 3000
	}
	portStr := strconv.Itoa(port)

	// mux := http.NewServeMux()
	mux := httprouter.New()
	wrapper := &httpWrapper{router: mux}

	// mux.HandleFunc("GET /{$}", dummyHandler("root"))
	// mux.HandleFunc("GET /something", dummyHandler("something"))
	// mux.HandleFunc("GET /something/{id}", dummyHandler("something/{id}"))
	// mux.HandleFunc("GET /something/{$}", dummyHandler("something/{$}"))
	// mux.HandleFunc("POST /somethingelse", dummyHandler("somethingelse"))

	// router := routes.NewRouter(mux)
	// router.GET("/{$}", dummyHandler("root"))
	// router.GET("/something", dummyHandler("something"))
	// router.GET("/something/{id}", dummyHandler("something/{id}"))
	// router.GET("/something/{$}", dummyHandler("something/{$}"))
	// router.POST("/somethingelse", dummyHandler("somethingelse"))

	mux.GET("/", wrapStdHandler(dummyHandler("root")))
	wrapper.GETWrapper("/something", dummyHandler("something"))
	wrapper.GETWrapper("/something/{id}", dummyHandler("something/{id}"))
	wrapper.GETWrapper("/something/{id}/{name}", dummyHandler("something/{id}/{name}"))
	wrapper.GETWrapper("/something/{id}/{name}/{wow}", dummyHandler("something/{id}/{name}/{wow}"))
	mux.POST("/somethingelse", wrapStdHandler(dummyHandler("somethingelse")))

	log.Printf("Starting server on port %s (Dev Mode: %t)...", portStr, isDev)
	err := http.ListenAndServe(":"+portStr, mux)
	if err != nil {
		log.Fatal(err)
	}

}

type httpWrapper struct {
	router *httprouter.Router
}

var re = regexp.MustCompile(`\{(\w+)\}`)

func (mx *httpWrapper) GETWrapper(path string, handler func(http.ResponseWriter, *http.Request)) {
	newpath := re.ReplaceAllString(path, `:$1`)
	mx.router.GET(newpath, wrapStdHandler(handler))
}

func dummyHandler(input string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		params := r.URL.Query()
		id := r.PathValue("id")
		fmt.Fprintf(w, "Path: %s\nParams: %v\nID: %s \nInput: %s", path, params, id, input)
	}
}

func dummyHandlerRouter(input string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r.SetPathValue("id", p.ByName("id"))
		fmt.Fprintf(w, "Path: %s\nParams: %v\nID: %v\nInput: %s %v\n Std Params:", r.URL.Path, r.URL.Query(), p.ByName("id"), input, r.PathValue("wow"))
	}
}

func wrapStdHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("wrapStdHandler", p)
		for _, param := range p {
			r.SetPathValue(param.Key, param.Value)
		}
		h.ServeHTTP(w, r)
	}
}
