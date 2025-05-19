package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRouteRegistrationAndDispatch tests basic routing
func TestRouteRegistrationAndDispatch(t *testing.T) {
	mux := http.NewServeMux()
	router := NewRouter(mux)

	// Dummy handler that sets a flag
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	router.ROUTE("/test-basic").GET(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test-basic", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if !handlerCalled {
		t.Error("Expected handler to be called, but it wasn't")
	}
}

// TestUnhandledMethod tests routing for unhandled methods
func TestUnhandledMethod(t *testing.T) {
	mux := http.NewServeMux()
	router := NewRouter(mux)

	// Register only GET
	router.ROUTE("/test-unhandled").GET(dummyHandlerGET)

	// Try a POST request
	req := httptest.NewRequest(http.MethodPost, "/test-unhandled", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rw.Code)
	}
}

// TestNonExistentRoute tests routing for paths not registered
func TestNonExistentRoute(t *testing.T) {
	mux := http.NewServeMux()
	router := NewRouter(mux)

	// Register a route, but test a different path
	router.ROUTE("/test-exists").GET(dummyHandlerGET)

	req := httptest.NewRequest(http.MethodGet, "/test-nonexistent", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rw.Code)
	}
}

// Add tests for Path Parameters here if you implement that functionality...
