package main

import (
	"log"
	"time"

	"mc-mono/go-server/database"
	"mc-mono/go-server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
	"github.com/shareed2k/goth_fiber"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	database.ConnectDB()

	app := fiber.New()

	// Session configuration
	store := session.New(session.Config{
		KeyLookup:      "cookie:session",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour,
		KeyGenerator:   utils.UUID,
	})

	goth_fiber.SessionStore = store

	// Setup routes
	router := routes.NewRouter(app)
	router.SetupRoutes()

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
