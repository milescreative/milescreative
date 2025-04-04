package routes

import (
	"mc-mono/go-server/internal/middleware"
	"mc-mono/go-server/routes/admin"
	"mc-mono/go-server/routes/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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
	// CSRF protection for all routes
	r.app.Use(middleware.SetupCSRF())

	// Session validation for protected routes
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
