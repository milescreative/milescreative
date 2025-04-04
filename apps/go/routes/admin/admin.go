package admin

import (
	"github.com/gofiber/fiber/v2"
)

type AdminRouter struct {
	app *fiber.App
}

func NewAdminRouter(app *fiber.App) *AdminRouter {
	return &AdminRouter{
		app: app,
	}
}

func (r *AdminRouter) SetupRoutes() {
	admin := r.app.Group("/admin")

	// Dashboard route
	admin.Get("/dashboard", r.handleDashboard)
}

func (r *AdminRouter) handleDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "This is a protected route",
	})
}
