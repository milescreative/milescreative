package main

import (
	"go-std/internal/auth"
	"go-std/internal/config"

	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go-std/internal/middleware"
	"go-std/internal/sqlc"
	"go-std/internal/utils"
	"go-std/routes"

	"github.com/g-h-miles/httpmux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	env, _ := config.Config()

	log.Println("env db url", env.GetString("DATABASE_URL"))

	dbpool, err := pgxpool.New(context.Background(), env.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	queries := sqlc.New(dbpool)

	now, err := queries.TestDatabaseConnection(context.Background())
	if err != nil || now == nil {
		log.Fatalf("Failed to test database connection: %v\n", err)
	}

	fmt.Printf("Successfully connected to database! Current time: %v\n", now)

	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	app := &config.App{
		DB:      dbpool,
		IsDev:   isDev,
		Queries: queries,
		Env:     env,
	}

	mux := httpmux.NewServeMux()
	// mux := http.NewServeMux()
	routes.AuthRoutes(mux, app)

	middlewareStack := middleware.CreateStack(middleware.CORS, middleware.CSPMiddleware(isDev))

	middlewareContext := middleware.NewMiddlewareContext(app)
	protected := middlewareContext.Protected

	authHandlers, err := auth.NewAuthHandlers(app)
	if err != nil {
		log.Fatalf("failed to create auth handlers: %v", err)
	}

	mux.HandleFunc("GET", "/api/auth/protected", protected(someProtectedHandler))

	mux.HandleFunc("GET", "/api/auth/test-form", middlewareStack(middleware.CSRFMiddleware(authHandlers.TestFormHandler)))
	mux.HandleFunc("GET", "/api/auth/csrf-protected", middlewareStack(middleware.CSRFMiddleware((someCSRFHandler))))
	mux.HandleFunc("GET", "/testing/{wow...}", someProtectedHandler)

	withMethodMiddleware := protectedMiddleware(mux, []string{"POST"}, testLogMiddleware)

	port, _ := env.GetInt("APP_PORT")
	if port == 0 {
		port = 3000
	}
	portStr := strconv.Itoa(port)

	log.Printf("Starting server on port %s (Dev Mode: %t)...", portStr, isDev)
	err = http.ListenAndServe(":"+portStr, middlewareStack(withMethodMiddleware.ServeHTTP))
	if err != nil {
		log.Fatal(err)
	}

}

func someProtectedHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, "You are authorized")
}

func someCSRFHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, "You are authorized")
}

func protectedMiddleware(router http.Handler, methods []string, mw func(http.Handler) http.Handler) http.Handler {
	methodSet := make(map[string]bool)
	for _, m := range methods {
		methodSet[m] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if methodSet[r.Method] {
			mw(router).ServeHTTP(w, r)
			return
		}
		router.ServeHTTP(w, r)
	})
}

func testLogMiddleware(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("testLogMiddleware")
		router.ServeHTTP(w, r)
	})
}
