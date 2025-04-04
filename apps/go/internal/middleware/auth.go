package middleware

import (
	"fmt"
	"mc-mono/go-server/internal/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected routes that require authentication
var protectedRoutes = map[string]bool{
	"/api/user/profile":  true,
	"/api/user/sessions": true,
	// Add more protected routes here
}

// Protected prefixes that require authentication
var protectedPrefixes = []string{
	"/api/user/",
	"/api/admin/",
	"/admin/",
	// Add more protected prefixes here
}

// isProtectedRoute checks if a route requires authentication
func isProtectedRoute(path string) bool {
	// Check exact matches
	if protectedRoutes[path] {
		return true
	}

	// Check prefixes
	for _, prefix := range protectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// RequireSession middleware checks for valid session except for public paths
func RequireSession() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := ctx.Path()
		method := ctx.Method()
		ip := ctx.IP()

		// Log every request
		fmt.Printf("\n🌐 [%s] %s - IP: %s\n", method, path, ip)

		// Check if route needs protection
		needsAuth := isProtectedRoute(path)
		if !needsAuth {
			fmt.Printf("✨ Public access: %s\n", path)
			return ctx.Next()
		}

		fmt.Printf("🔒 Checking authentication for: %s\n", path)

		// Validate session for all other routes
		sessionData, err := auth.ValidateSession(ctx)
		if err != nil {
			fmt.Printf("❌ Auth failed: %s - Error: %s\n", path, err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"authenticated": false,
				"message":       err.Error(),
			})
		}

		// Log successful auth
		fmt.Printf("✅ Authenticated user: %s (%s)\n", sessionData.User.Name, sessionData.User.Email)

		// Add session data to context
		ctx.Locals("session", sessionData.Session)
		ctx.Locals("user", sessionData.User)

		return ctx.Next()
	}
}
