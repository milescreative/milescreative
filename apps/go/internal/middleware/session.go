package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store

func SetupSession() *session.Store {
	sessConfig := session.Config{
		Expiration:     30 * time.Minute,        // Expire sessions after 30 minutes of inactivity
		KeyLookup:      "cookie:__Host-session", // Recommended to use the __Host- prefix when serving the app over TLS
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	}
	SessionStore := session.New(sessConfig)

	return SessionStore
}
