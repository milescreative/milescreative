package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"mc-mono/go-server/database"
	"mc-mono/go-server/database/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionData struct {
	Session models.Session
	User    models.User
}

const (
	SessionCookieName = "session"
	SessionDuration   = 24 * time.Hour
)

func CreateSessionCookie(token string, expiresAt time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
}

func ClearSessionCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
}

func CreateSession(userUID string, ctx *fiber.Ctx) (*models.Session, string, error) {
	sessionToken, err := GenerateSecureToken(32)
	if err != nil {
		return nil, "", err
	}

	session := models.Session{
		UserUID:   userUID,
		UserAgent: ctx.Get("User-Agent"),
		IPAddress: ctx.IP(),
		ExpiresAt: time.Now().Add(SessionDuration),
	}
	session.SetToken(sessionToken)

	if err := database.DB.Create(&session).Error; err != nil {
		return nil, "", err
	}

	return &session, sessionToken, nil
}

func FormatUserResponse(user models.User) fiber.Map {
	return fiber.Map{
		"uid":      user.UID,
		"email":    user.Email,
		"name":     user.Name,
		"avatar":   user.AvatarURL,
		"provider": user.Provider,
	}
}

// ValidateSession checks for valid session and returns session data
func ValidateSession(ctx *fiber.Ctx) (*SessionData, error) {
	sessionToken := ctx.Cookies("session")
	if sessionToken == "" {
		fmt.Printf("🔑 No session cookie present\n")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "No active session")
	}

	var session models.Session
	tokenHash := models.HashToken(sessionToken)

	// Find active session
	err := database.DB.Where("token_hash = ? AND expires_at > ?",
		tokenHash,
		time.Now(),
	).Preload("User").First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("🔑 Session check: No active session found for token\n")
		} else {
			fmt.Printf("❌ Session check error: %v\n", err)
		}

		// Clear invalid session cookie
		ctx.Cookie(ClearSessionCookie())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired session")
	}

	fmt.Printf("✅ Valid session found for user: %s\n", session.User.Email)

	// Update last used time
	database.DB.Save(&session)

	return &SessionData{
		Session: session,
		User:    session.User,
	}, nil
}

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
