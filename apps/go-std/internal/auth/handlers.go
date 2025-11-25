package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"go-std/internal/config"
	"go-std/internal/utils"
	"net/http"

	"strings"
	"time"

	"go-std/internal/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

var logger = utils.NewLogger(utils.DEBUG, true)

type AuthHandlers struct {
	*config.App
	SessionCookieName           string
	SessionExpiration           time.Duration
	AuthRedirectQueryParam      string
	AuthRedirectCookieName      string
	AuthRedirectDefault         string
	OAuthStateCookieName        string
	OAuthCodeVerifierCookieName string
	UserSessionQueryParam       string
	ProviderRegistry            *ProviderRegistry
}

const (
	sessionCookieName           = "session_token"
	sessionExpiration           = time.Hour * 24 * 30
	authRedirectQueryParam      = "redirect_url"
	authRedirectCookieName      = "auth_redirect_url"
	authRedirectDefault         = "/"
	oauthStateCookieName        = "oauth_state"
	oauthCodeVerifierCookieName = "oauth_code_verifier"
	userSessionQueryParam       = "user_session"
)

func NewAuthHandlers(app *config.App) (*AuthHandlers, error) {
	registry, err := NewProviderRegistry(app)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider registry: %w", err)
	}

	return &AuthHandlers{
		App:                         app,
		SessionCookieName:           sessionCookieName,
		SessionExpiration:           sessionExpiration,
		AuthRedirectQueryParam:      authRedirectQueryParam,
		AuthRedirectCookieName:      authRedirectCookieName,
		AuthRedirectDefault:         authRedirectDefault,
		OAuthStateCookieName:        oauthStateCookieName,
		OAuthCodeVerifierCookieName: oauthCodeVerifierCookieName,
		UserSessionQueryParam:       userSessionQueryParam,
		ProviderRegistry:            registry,
	}, nil
}

func (a *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {

	provider := r.PathValue("provider")
	if provider == "" {
		provider = "google" // default provider
	}

	q := a.Queries
	redirectURL := r.URL.Query().Get(a.AuthRedirectQueryParam)
	if redirectURL == "" {
		redirectURL = a.AuthRedirectDefault
	}

	// check if user is already logged in
	valid_session, _ := utils.ValidateSession(q, w, r)
	if valid_session {
		utils.Redirect(w, r, redirectURL)
		return
	}

	if redirectURL != "" {
		cookie := &http.Cookie{
			Name:   a.AuthRedirectCookieName,
			Value:  redirectURL,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}

	state, err := utils.GenerateState()
	if err != nil {
		logger.Error("Error generating state: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error generating state", "INTERNAL_SERVER_ERROR")
		return
	}
	// logger.Debug("state: %s", state)

	codeVerifier, err := utils.GenerateCodeVerifier()
	if err != nil {
		logger.Error("Error generating codeVerifier: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error generating codeVerifier", "INTERNAL_SERVER_ERROR")
		return
	}
	// logger.Debug("codeVerifier: %s", codeVerifier)

	// Create OAuth provider
	oauthProvider, err := a.ProviderRegistry.CreateProvider(provider)
	if err != nil {
		logger.Error("Error creating OAuth provider: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Unsupported provider", "BAD_REQUEST")
		return
	}

	authURL, err := oauthProvider.CreateAuthorizationURL(state, codeVerifier)
	if err != nil {
		logger.Error("Error creating authorization URL: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating authorization URL", "INTERNAL_SERVER_ERROR")
		return
	}

	isDev := a.IsDev

	cookies := []*http.Cookie{{
		Name:     a.OAuthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
	},
		{
			Name:     a.OAuthCodeVerifierCookieName,
			Value:    codeVerifier,
			Path:     "/",
			HttpOnly: true,
			Secure:   !isDev,
			SameSite: http.SameSiteLaxMode,
		},
		{
			Name:     "oauth_provider",
			Value:    provider,
			Path:     "/",
			HttpOnly: true,
			Secure:   !isDev,
			SameSite: http.SameSiteLaxMode,
		},
	}
	for _, c := range cookies {
		http.SetCookie(w, c)
	}

	// logger.Debug("login creation success- redirecting to: %s", authURL)
	// utils.SuccessResponse(w, u)
	http.Redirect(w, r, authURL.String(), http.StatusSeeOther)

}

func (a *AuthHandlers) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	q := a.Queries
	storedState, err := r.Cookie(a.OAuthStateCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting oauth_state. Please restart", "BAD_REQUEST")
		return
	}
	storedCodeVerifier, err := r.Cookie(a.OAuthCodeVerifierCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting oauth_code_verifier. Please restart", "BAD_REQUEST")
		return
	}

	provider := r.PathValue("provider")
	if provider == "" {
		provider = r.URL.Query().Get("provider")
	}
	if provider == "" {
		providerCookie, err := r.Cookie("oauth_provider")
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Error getting oauth_provider. Please restart", "BAD_REQUEST")
			return
		}
		provider = providerCookie.Value
	}

	if provider == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting provider. Please restart", "BAD_REQUEST")
		return
	}

	if storedState.Value != state {
		utils.ErrorResponse(w, http.StatusBadRequest, "state mismatch- please restart", "BAD_REQUEST")
		return
	}
	// Create OAuth provider
	oauthProvider, err := a.ProviderRegistry.CreateProvider(provider)
	if err != nil {
		logger.Error("Error creating OAuth provider: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Unsupported provider", "BAD_REQUEST")
		return
	}

	tokens, err := oauthProvider.ValidateAuthorizationCode(code, storedCodeVerifier.Value)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error validating authorization code. Please restart", "BAD_REQUEST")
		logger.Error("Error validating authorization code: %v", err)
		return
	}

	// Clean up cookies
	utils.RemoveCookie(w, a.OAuthStateCookieName)
	utils.RemoveCookie(w, a.OAuthCodeVerifierCookieName)
	utils.RemoveCookie(w, "oauth_provider")

	// Get user info from provider
	userInfo, err := oauthProvider.GetUserInfo(tokens)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting user info", "BAD_REQUEST")
		logger.Error("Error getting user info: %v", err)
		return
	}

	token_result, err := tokens.GetTokenResult()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting token result", "BAD_REQUEST")
		logger.Error("Error getting token result: %v", err)
		return
	}

	ip_address := r.RemoteAddr
	user_agent := r.UserAgent()

	session_token, err := utils.GenerateSessionToken()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error generating session token", "BAD_REQUEST")
		logger.Error("Error generating session token: %v", err)
		return
	}

	_, thisErr := q.CreateNewUser(context.Background(), sqlc.CreateNewUserParams{
		Name:                 userInfo.Name,
		Email:                userInfo.Email,
		EmailVerified:        userInfo.EmailVerified,
		Image:                pgtype.Text{String: userInfo.Picture, Valid: true},
		AccountID:            userInfo.ID,
		ProviderID:           oauthProvider.GetProviderName(),
		Scope:                pgtype.Text{String: strings.Join(token_result.Scopes, " "), Valid: true},
		AccessToken:          pgtype.Text{String: token_result.AccessToken, Valid: true},
		RefreshToken:         pgtype.Text{String: token_result.RefreshToken, Valid: true},
		IDToken:              pgtype.Text{String: token_result.IDToken, Valid: true},
		ExpiresAt:            pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 30), Valid: true},
		Token:                session_token,
		IpAddress:            pgtype.Text{String: ip_address, Valid: true},
		UserAgent:            pgtype.Text{String: user_agent, Valid: true},
		AccessTokenExpiresAt: pgtype.Timestamp{Time: token_result.AccessTokenExpiresAt, Valid: true},
	})

	if thisErr != nil {
		logger.Error("Error creating new user: %v", thisErr)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating new user", "INTERNAL_SERVER_ERROR")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   session_token,
		Path:    "/",
		Expires: time.Now().Add(sessionExpiration),
	})

	redirectURL, err := r.Cookie(authRedirectCookieName)
	if err != nil {
		logger.Error("Error getting auth_redirect_url: %v", err)
	}

	if redirectURL != nil && redirectURL.Value != "" {
		utils.RemoveCookie(w, authRedirectCookieName)
		utils.Redirect(w, r, redirectURL.Value)
		return
	} else {
		utils.Redirect(w, r, authRedirectDefault)
		return
	}

}

func (a *AuthHandlers) ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ValidateSessionHandler")
	q := a.Queries

	valid_session, err := utils.ValidateSession(q, w, r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error validating session", "INTERNAL_SERVER_ERROR")
		logger.Error("error validating session: %v", err)
		return
	}
	if valid_session {
		utils.SuccessResponse(w, "session valid")
		return
	}

	utils.ErrorResponse(w, http.StatusUnauthorized, "session invalid", "UNAUTHORIZED")

}

func (a *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	q := a.Queries

	cookie, _ := r.Cookie(sessionCookieName)
	if cookie != nil {
		_, err := q.DeleteSession(context.Background(), cookie.Value)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error logging out", "INTERNAL_SERVER_ERROR")
			logger.Error("error logging out: %v", err)
			return
		}
		utils.RemoveCookie(w, sessionCookieName)
	}

	utils.SuccessResponse(w, "session logged out")
}

func (a *AuthHandlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	q := a.Queries

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error getting session cookie", "BAD_REQUEST")
		logger.Error("error getting session cookie: %v", err)
		return
	}

	session, err := q.GetSessionByToken(context.Background(), cookie.Value)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error getting session", "BAD_REQUEST")
		logger.Error("error getting session: %v", err)
		return
	}

	// Get the provider from the session
	providerName := session.ProviderID
	if providerName == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "no provider found for session", "BAD_REQUEST")
		logger.Error("no provider found for session")
		return
	}

	// Create OAuth provider
	oauthProvider, err := a.ProviderRegistry.CreateProvider(providerName)
	if err != nil {
		logger.Error("Error creating OAuth provider: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Unsupported provider", "BAD_REQUEST")
		return
	}

	// Check if provider supports refresh tokens
	refreshToken := session.RefreshToken.String
	if refreshToken == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "no refresh token available", "BAD_REQUEST")
		logger.Error("no refresh token available for session")
		return
	}

	tokens, err := oauthProvider.RefreshAccessToken(refreshToken)
	if err != nil {
		// Handle provider-specific errors (like GitHub not supporting refresh)
		if providerName == "github" {
			utils.ErrorResponse(w, http.StatusBadRequest, "GitHub tokens do not expire and cannot be refreshed", "BAD_REQUEST")
		} else {
			utils.ErrorResponse(w, http.StatusBadRequest, "error refreshing access token", "BAD_REQUEST")
		}
		logger.Error("error refreshing access token: %v", err)
		return
	}

	accessToken, _ := tokens.AccessToken()
	idToken, _ := tokens.IDToken()
	expiresAt, _ := tokens.AccessTokenExpiresAt()

	err = q.UpdateAccount(context.Background(), sqlc.UpdateAccountParams{
		AccountID:            session.AccountID.String,
		AccessToken:          pgtype.Text{String: accessToken, Valid: true},
		IDToken:              pgtype.Text{String: idToken, Valid: true},
		AccessTokenExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error updating account", "INTERNAL_SERVER_ERROR")
		logger.Error("error updating account: %v", err)
		return
	}

	utils.SuccessResponse(w, "refresh token successful")
}

func (a *AuthHandlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	q := a.Queries

	user_id := r.URL.Query().Get(userSessionQueryParam)
	if user_id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user_id is required", "BAD_REQUEST")
		return
	}
	user, err := q.GetUserByID(context.Background(), user_id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error getting user", "INTERNAL_SERVER_ERROR")
		logger.Error("error getting user: %v", err)
		return
	}
	utils.SuccessResponse(w, user)
}

func (a *AuthHandlers) GetUserSessionsHandler(w http.ResponseWriter, r *http.Request) {
	q := a.Queries

	user_id := r.URL.Query().Get(userSessionQueryParam)
	sessions, err := q.GetUserSessions(context.Background(), user_id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error getting user sessions", "INTERNAL_SERVER_ERROR")
		logger.Error("error getting user sessions: %v", err)
		return
	}
	utils.SuccessResponse(w, sessions)
}

func (a *AuthHandlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	q := a.Queries

	user_id := r.URL.Query().Get(userSessionQueryParam)
	if user_id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user_id is required", "BAD_REQUEST")
		return
	}
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	image := r.URL.Query().Get("image")

	id, err := q.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		ID:    user_id,
		Name:  name,
		Email: email,
		Image: pgtype.Text{String: image, Valid: true},
	})
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user not found", "BAD_REQUEST")
		return
	}
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error updating user", "INTERNAL_SERVER_ERROR")
		logger.Error("error updating user: %v", err)
		return
	}
	utils.SuccessResponse(w, "user updated")
}

func (a *AuthHandlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	q := a.Queries

	user_id := r.URL.Query().Get(userSessionQueryParam)
	id, err := q.DeleteUser(context.Background(), user_id)
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user not found", "BAD_REQUEST")
		return
	}
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error deleting user", "INTERNAL_SERVER_ERROR")
		logger.Error("error deleting user: %v", err)
		return
	}
	utils.SuccessResponse(w, "user deleted")
}

func (a *AuthHandlers) GetCSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Session required", "UNAUTHORIZED")
		return
	}

	utils.SetCSRFToken(w, sessionCookie.Value)
	utils.SuccessResponse(w, "CSRF token set")
}

func (a *AuthHandlers) TestFormHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	// Parse the JSON body
	var formData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", "BAD_REQUEST")
		return
	}

	// Log the form data
	logger.Debug("Received form data: %+v", formData)

	// Return the form data in the response
	utils.SuccessResponse(w, formData)
}
