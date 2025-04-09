package auth

import (
	"mc-mono/go-server/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

type AuthRouter struct {
	app      *fiber.App
	handlers *auth.AuthHandlers
}

func NewAuthRouter(app *fiber.App) *AuthRouter {
	return &AuthRouter{
		app:      app,
		handlers: auth.NewAuthHandlers(),
	}
}

func (r *AuthRouter) SetupRoutes() {

	// Configure Oauth
	SetupProviders()

	// Auth routes
	r.app.Post("/api/auth/login/:provider", goth_fiber.BeginAuthHandler)
	r.app.Get("/api/auth/callback/:provider", r.handlers.HandleGoogleCallback)
	r.app.Get("/api/auth/success", r.handlers.HandleSuccess)
	r.app.Post("/api/auth/logout", r.handlers.HandleLogout)
	r.app.Get("/api/auth/csrf-token", r.handlers.HandleCSRFToken)
	r.app.Get("/api/auth/status", r.handlers.HandleStatus)
	r.app.Get("/api/auth/sessions", r.handlers.HandleListSessions)
	r.app.Post("/api/auth/sessionsPost", r.handlers.HandleListSessions)
	r.app.Post("api/auth/sessions/delete", r.handlers.HandleDeleteSessions)
}
