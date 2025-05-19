package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// dummyHandlerGET, dummyHandlerPOST, dummyHandlerDELETE are distinct handlers
func dummyHandlerGET(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func dummyHandlerPOST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // Different status to distinguish
}

func dummyHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent) // Different status to distinguish
}

// BenchmarkDirectServeMuxMultipleHandlers benchmarks using direct mux.HandleFunc with distinct handlers
func BenchmarkDirectServeMuxMultipleHandlers(b *testing.B) {
	mux := http.NewServeMux()

	// Register distinct handlers directly with methods
	mux.HandleFunc(http.MethodGet+"/benchmark-distinct", dummyHandlerGET)
	mux.HandleFunc(http.MethodPost+"/benchmark-distinct", dummyHandlerPOST)
	mux.HandleFunc(http.MethodDelete+"/benchmark-distinct", dummyHandlerDELETE)
	mux.HandleFunc("/benchmark-nondistinct", dummyHandlerGET)

	reqGET := httptest.NewRequest(http.MethodGet, "/benchmark-distinct", nil)
	reqPOST := httptest.NewRequest(http.MethodPost, "/benchmark-distinct", nil)
	reqDELETE := httptest.NewRequest(http.MethodDelete, "/benchmark-distinct", nil)
	reqNONDISTINCT := httptest.NewRequest(http.MethodGet, "/benchmark-nondistinct", nil)
	rw := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, reqGET)
		mux.ServeHTTP(rw, reqPOST)
		mux.ServeHTTP(rw, reqDELETE)
		mux.ServeHTTP(rw, reqNONDISTINCT)
	}
}

// BenchmarkChainedRouterMultipleHandlers benchmarks using the chained router syntax with distinct handlers
func BenchmarkChainedRouterMultipleHandlers(b *testing.B) {
	mux := http.NewServeMux()
	router := NewRouter(mux)

	// Register distinct handlers using the chained syntax
	router.ROUTE("/benchmark-distinct").
		GET(dummyHandlerGET).
		POST(dummyHandlerPOST).
		DELETE(dummyHandlerDELETE)

	router.ROUTE("/benchmark-nondistinct").
		GET(dummyHandlerGET)

	reqGET := httptest.NewRequest(http.MethodGet, "/benchmark-distinct", nil)
	reqPOST := httptest.NewRequest(http.MethodPost, "/benchmark-distinct", nil)
	reqDELETE := httptest.NewRequest(http.MethodDelete, "/benchmark-distinct", nil)
	reqNONDISTINCT := httptest.NewRequest(http.MethodGet, "/benchmark-nondistinct", nil)
	rw := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, reqGET)
		mux.ServeHTTP(rw, reqPOST)
		mux.ServeHTTP(rw, reqDELETE)
		mux.ServeHTTP(rw, reqNONDISTINCT)
	}
}
