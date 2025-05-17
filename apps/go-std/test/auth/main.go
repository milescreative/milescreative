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

	middlewareStack := middleware.CreateStack(middleware.CORS, middleware.CSPMiddleware(isDev))

	middlewareContext := middleware.NewMiddlewareContext(app)
	protected := middlewareContext.Protected
	authHandlers := auth.NewAuthHandlers(app)

	mux.HandleFunc("/", authHandlers.EnvHandler)
	mux.HandleFunc("/api/auth/login", authHandlers.LoginHandler)
	mux.HandleFunc("/api/auth/callback/google", authHandlers.CallbackHandler)
	mux.HandleFunc("/api/auth/validate", authHandlers.ValidateSessionHandler)
	mux.HandleFunc("/api/auth/logout", authHandlers.LogoutHandler)
	mux.HandleFunc("/api/auth/refresh", authHandlers.RefreshTokenHandler)

	mux.HandleFunc("/api/auth/protected", protected(someProtectedHandler))

	mux.HandleFunc("/api/auth/csrf", authHandlers.GetCSRFTokenHandler)
	mux.Handle("/api/auth/test-form", middlewareStack(middleware.CSRFMiddleware(http.HandlerFunc(authHandlers.TestFormHandler))))
	mux.Handle("/api/auth/csrf-protected", middlewareStack(middleware.CSRFMiddleware(http.HandlerFunc(someCSRFHandler))))

	port, _ := env.GetInt("APP_PORT")
	if port == 0 {
		port = 3000
	}
	portStr := strconv.Itoa(port)

	log.Printf("Starting server on port %s (Dev Mode: %t)...", portStr, isDev)
	err = http.ListenAndServe(":"+portStr, middlewareStack(mux))
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
