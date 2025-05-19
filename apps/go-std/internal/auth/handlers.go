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
}

const (
	sessionCookieName            = "session_token"
	sessionExpiration            = time.Hour * 24 * 30
	authRedirectQueryParam       = "redirect_url"
	authRedirectCookieName       = "auth_redirect_url"
	authRedirectDefault          = "/"
	googleOAuthStateCookieName   = "google_oauth_state"
	googleCodeVerifierCookieName = "google_code_verifier"
	userSessionQueryParam        = "user_session"
)

func NewAuthHandlers(app *config.App) *AuthHandlers {
	return &AuthHandlers{
		App: app,
	}
}

func (a *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {

	q := a.Queries
	redirectURL := r.URL.Query().Get(authRedirectQueryParam)
	if redirectURL == "" {
		redirectURL = authRedirectDefault
	}

	// check if user is already logged in
	valid_session, _ := utils.ValidateSession(q, w, r)
	if valid_session {
		utils.Redirect(w, r, redirectURL)
		return
	}

	if redirectURL != "" {
		cookie := &http.Cookie{
			Name:   authRedirectCookieName,
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
	logger.Debug("state: %s", state)

	codeVerifier, err := utils.GenerateCodeVerifier()
	if err != nil {
		logger.Error("Error generating codeVerifier: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error generating codeVerifier", "INTERNAL_SERVER_ERROR")
		return
	}
	logger.Debug("codeVerifier: %s", codeVerifier)

	isDev := a.IsDev

	goog := NewGoogleOAuth(
		a.Env.GetString("GOOGLE_CLIENT_ID"),
		a.Env.GetString("GOOGLE_CLIENT_SECRET"),
		a.Env.GetString("GOOGLE_REDIRECT_URI"),
		[]string{"email", "profile"},
	)

	u, err := goog.CreateAuthorizationURLWithPKCE(state, codeVerifier)
	if err != nil {
		logger.Error("Error creating authorization URL: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating authorization URL", "INTERNAL_SERVER_ERROR")
		return
	}
	logger.Debug("u: %s", u)

	cookies := []*http.Cookie{{
		Name:     googleOAuthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
	},
		{
			Name:     googleCodeVerifierCookieName,
			Value:    codeVerifier,
			Path:     "/",
			HttpOnly: true,
			Secure:   !isDev,
			SameSite: http.SameSiteLaxMode,
		},
	}
	for _, c := range cookies {
		http.SetCookie(w, c)
	}

	logger.Debug("login creation success- redirecting to: %s", u)
	// utils.SuccessResponse(w, u)
	http.Redirect(w, r, u.String(), http.StatusSeeOther)

}

func (a *AuthHandlers) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	q := a.Queries
	storedState, err := r.Cookie(googleOAuthStateCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting google_oauth_state. Please restart", "BAD_REQUEST")
		return
	}
	storedCodeVerifier, err := r.Cookie(googleCodeVerifierCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting google_code_verifier. Please restart", "BAD_REQUEST")
		return
	}

	if storedState.Value != state {
		utils.ErrorResponse(w, http.StatusBadRequest, "state mismatch- please restart", "BAD_REQUEST")
		return
	}
	goog := NewGoogleOAuth(
		a.Env.GetString("GOOGLE_CLIENT_ID"),
		a.Env.GetString("GOOGLE_CLIENT_SECRET"),
		a.Env.GetString("GOOGLE_REDIRECT_URI"),
		[]string{"email", "profile"},
	)
	tokens, err := goog.ValidateAuthorizationCode(code, storedCodeVerifier.Value)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error validating authorization code. Please restart", "BAD_REQUEST")
		logger.Error("Error validating authorization code: %v", err)
		return
	}
	utils.RemoveCookie(w, googleOAuthStateCookieName)
	utils.RemoveCookie(w, googleCodeVerifierCookieName)

	token_result, err := tokens.GetTokenResult()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error getting token result", "BAD_REQUEST")
		logger.Error("Error getting token result: %v", err)
		return
	}

	claims, err := utils.DecodeJwt(token_result.IDToken)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Error decoding idToken", "BAD_REQUEST")
		logger.Error("Error decoding idToken: %v", err)
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
		Name:                 claims["name"].(string),
		Email:                claims["email"].(string),
		EmailVerified:        claims["email_verified"].(bool),
		Image:                pgtype.Text{String: claims["picture"].(string), Valid: true},
		AccountID:            claims["sub"].(string),
		ProviderID:           "google",
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

// TODO: remove this
func (a *AuthHandlers) EnvHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	str := fmt.Sprintf(`
        <h1>Environment Variables</h1>
        <pre>
APP_ENV: %s
GOOGLE_CLIENT_ID: %s
GOOGLE_CLIENT_SECRET: %s
GOOGLE_REDIRECT_URI: %s
PORT: %d
    `, a.Env.GetString("APP_ENV"),
		a.Env.GetString("GOOGLE_CLIENT_ID"),
		a.Env.GetString("GOOGLE_CLIENT_SECRET"),
		a.Env.GetString("GOOGLE_REDIRECT_URI"),
		a.Env.Port(),
	)
	utils.SuccessResponse(w, str)
}

func (a *AuthHandlers) ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
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
	utils.PostOrBail(w, r)

	q := a.Queries

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error getting session cookie", "BAD_REQUEST")
		logger.Error("error getting session cookie: %v", err)
		return
	}

	goog := NewGoogleOAuth(
		a.Env.GetString("GOOGLE_CLIENT_ID"),
		a.Env.GetString("GOOGLE_CLIENT_SECRET"),
		a.Env.GetString("GOOGLE_REDIRECT_URI"),
		[]string{"email", "profile"},
	)

	session, err := q.GetSessionByToken(context.Background(), cookie.Value)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error getting session", "BAD_REQUEST")
		logger.Error("error getting session: %v", err)
		return
	}
	refresh_token := session.RefreshToken.String
	tokens, err := goog.RefreshAccessToken(refresh_token)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error refreshing access token", "BAD_REQUEST")
		logger.Error("error refreshing access token: %v", err)
		return
	}

	access_token, _ := tokens.AccessToken()
	id_token, _ := tokens.IDToken()
	expires_at, _ := tokens.AccessTokenExpiresAt()

	err = q.UpdateAccount(context.Background(), sqlc.UpdateAccountParams{
		AccountID:            session.AccountID.String,
		AccessToken:          pgtype.Text{String: access_token, Valid: true},
		IDToken:              pgtype.Text{String: id_token, Valid: true},
		AccessTokenExpiresAt: pgtype.Timestamp{Time: expires_at, Valid: true},
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
