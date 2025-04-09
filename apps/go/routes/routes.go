package routes

import (
	"mc-mono/go-server/internal/middleware"
	"mc-mono/go-server/routes/admin"
	"mc-mono/go-server/routes/auth"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Router struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) *Router {
	return &Router{
		app: app,
	}
}

func (r *Router) setupMiddleware() {
	r.app.Use(favicon.New(favicon.Config{
		File: "./public/favicon.ico",
		URL:  "/favicon.ico",
	}))

	// Initialize a session store
	sessConfig := session.Config{
		Expiration: 30 * time.Minute, // Expire sessions after 30 minutes of inactivity
		// KeyLookup:      "cookie:__Host-session", // Recommended to use the __Host- prefix when serving the app over TLS
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	}
	store := session.New(sessConfig)
	csrfConfig := csrf.Config{
		Session: store,
		// KeyLookup:      "header:X-CSRF-Token", // In this example, we will be using a hidden input field to store the CSRF token
		CookieName:     "csrf_",
		CookieSameSite: "None", // Recommended to set this to Lax or Strict
		CookieSecure:   false,  // Recommended to set to true when serving the app over TLS
		CookieHTTPOnly: false,  // Recommended, otherwise if using JS framework recomend: false and KeyLookup: "header:X-CSRF-Token"
		ContextKey:     "csrf",
		Expiration:     30 * time.Minute,
		// CookieDomain:   "localhost:3002",
	}
	csrfMiddleware := csrf.New(csrfConfig)
	// middleware.SetupSession()
	// r.app.Use(middleware.SetupCSRF())
	r.app.Use(csrfMiddleware)
	r.app.Use(middleware.SetupCORS())
	r.app.Use(middleware.RequireSession())

	// Add future global middleware here
	// Example:
	// r.app.Use(middleware.Logger())
	// r.app.Use(middleware.RateLimit())
}

func (r *Router) SetupRoutes() {
	r.setupMiddleware()

	r.app.Static("/", "./public")
	// Setup auth routes
	authRouter := auth.NewAuthRouter(r.app)
	authRouter.SetupRoutes()

	adminRouter := admin.NewAdminRouter(r.app)
	adminRouter.SetupRoutes()

	// Add other route groups here as needed
	// Example:
	// apiRouter := api.NewApiRouter(r.app)
	// apiRouter.SetupRoutes()
}
