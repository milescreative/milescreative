package main

import (
	"go-std/internal/auth"
	"go-std/internal/config"

	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go-std/internal/middleware"
	"go-std/internal/sqlc"
	"go-std/internal/utils"
	"go-std/routes"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	env, _ := config.Config()
	dbpool, err := pgxpool.New(context.Background(), env.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	isDev := strings.ToLower(env.GetString("APP_ENV")) == "development"

	app := &config.App{
		DB:      dbpool,
		IsDev:   isDev,
		Queries: sqlc.New(dbpool),
		Env:     env,
	}

	mux := http.NewServeMux()
	routes.AuthRoutes(mux, app)

	middlewareStack := middleware.CreateStack(middleware.CORS, middleware.CSPMiddleware(isDev))

	middlewareContext := middleware.NewMiddlewareContext(app)
	protected := middlewareContext.Protected

	authHandlers := auth.NewAuthHandlers(app)

	mux.HandleFunc("/api/auth/protected", protected(someProtectedHandler))

	mux.HandleFunc("/api/auth/test-form", middlewareStack(middleware.CSRFMiddleware(authHandlers.TestFormHandler)))
	mux.HandleFunc("/api/auth/csrf-protected", middlewareStack(middleware.CSRFMiddleware((someCSRFHandler))))

	port, _ := env.GetInt("APP_PORT")
	if port == 0 {
		port = 3000
	}
	portStr := strconv.Itoa(port)

	log.Printf("Starting server on port %s (Dev Mode: %t)...", portStr, isDev)
	err = http.ListenAndServe(":"+portStr, middlewareStack(mux.ServeHTTP))
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
