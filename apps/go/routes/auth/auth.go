package auth

import (
	"os"
	"time"

	"mc-mono/go-server/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
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

	// middlewares
	r.app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}))

	// Configure Google OAuth2
	goth.UseProviders(
		google.New(os.Getenv("OAUTH_KEY"),
			os.Getenv("OAUTH_SECRET"),
			"http://localhost:3000/api/auth/callback/google",
			"email",
			"profile",
			"openid",
		))

	// Auth routes
	r.app.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.app.Get("/api/auth/callback/:provider", r.handlers.HandleGoogleCallback)
	r.app.Get("/auth/success", r.handlers.HandleSuccess)
	r.app.Get("/logout", r.handlers.HandleLogout)
	r.app.Get("/api/csrf-token", r.handlers.HandleCSRFToken)
	r.app.Get("/auth/status", r.handlers.HandleStatus)

}
