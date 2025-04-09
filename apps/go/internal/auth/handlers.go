package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"

	"mc-mono/go-server/database"
	"mc-mono/go-server/database/models"
)

type AuthHandlers struct{}

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

func (h *AuthHandlers) HandleGoogleCallback(ctx *fiber.Ctx) error {
	gothUser, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	// Check if user already exists
	var user models.User
	result := database.DB.Where("email = ?", gothUser.Email).First(&user)
	if result.Error != nil {
		// Create new user if not found
		user = models.User{
			Email:      gothUser.Email,
			Name:       gothUser.Name,
			Provider:   gothUser.Provider,
			ProviderID: gothUser.UserID,
			AvatarURL:  gothUser.AvatarURL,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	} else {
		// Update existing user
		user.LastLoginAt = time.Now()
		user.Name = gothUser.Name
		user.AvatarURL = gothUser.AvatarURL
		if err := database.DB.Save(&user).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}
	}

	session, sessionToken, err := CreateSession(user.UID, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	ctx.Cookie(CreateSessionCookie(sessionToken, session.ExpiresAt))

	return ctx.Redirect(fmt.Sprintf("%s/auth/success?uid=%s", os.Getenv("FRONTEND_URL"), user.UID))
}

func (h *AuthHandlers) HandleSuccess(ctx *fiber.Ctx) error {
	uid := ctx.Query("uid")
	if uid == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No UID provided",
		})
	}

	var user models.User
	if err := database.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"user":    FormatUserResponse(user),
	})
}

// HandleStatus checks the current auth status
func (h *AuthHandlers) HandleStatus(ctx *fiber.Ctx) error {
	sessionData, err := ValidateSession(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"authenticated": false,
			"message":       err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"authenticated": true,
		"user":          FormatUserResponse(sessionData.User),
		"session": fiber.Map{
			"expires_at": sessionData.Session.ExpiresAt,
			"user_agent": sessionData.Session.UserAgent,
		},
	})
}

// HandleLogout processes logout requests
func (h *AuthHandlers) HandleLogout(ctx *fiber.Ctx) error {
	sessionToken := ctx.Cookies(SessionCookieName)
	if sessionToken != "" {
		tokenHash := models.HashToken(sessionToken)
		database.DB.Where("token_hash = ?", tokenHash).Delete(&models.Session{})
	}

	ctx.Cookie(ClearSessionCookie())

	return ctx.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
func (h *AuthHandlers) HandleCSRFToken(ctx *fiber.Ctx) error {
	token := ctx.Locals("csrf")
	fmt.Printf("Generated CSRF token: %s\n", token)

	if token == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate CSRF token",
		})
	}

	return ctx.JSON(fiber.Map{
		"csrf_token": token,
	})
}

func (h *AuthHandlers) HandleListSessions(ctx *fiber.Ctx) error {
	sessionData, err := ValidateSession(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userUID := sessionData.User.UID
	sessions, err := GetSessions(userUID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get sessions",
		})
	}

	return ctx.JSON(fiber.Map{
		"sessions": sessions,
	})

}

func (h *AuthHandlers) HandleDeleteSessions(ctx *fiber.Ctx) error {
	sessionData, err := ValidateSession(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userUID := sessionData.User.UID
	err = DeleteUserSessions(ctx, userUID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sessions",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Sessions deleted successfully",
	})
}
