package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3002",
		AllowHeaders:     "*",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	})
}
