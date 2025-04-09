package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

const HeaderName = "X-CSRF-Token"

// SetupCSRF returns configured CSRF middleware
func SetupCSRF() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:" + HeaderName,
		CookieName:     "__Host-csrf",
		CookieSameSite: "Lax",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		Expiration:     30 * time.Minute,
		KeyGenerator:   utils.UUIDv4,
		Extractor:      csrf.CsrfFromHeader(HeaderName),
		TokenLookup:    "cookie:__Host-csrf",
		CookieDomain:   "",
		CookiePath:     "/",
		Session:        SessionStore,
		ContextKey:     "csrf",
	})
}
