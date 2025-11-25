package auth

import (
	"fmt"
	"go-std/internal/config"
	"go-std/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	googleAuthEndpoint   = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenEndpoint  = "https://oauth2.googleapis.com/token"
	googleRevokeEndpoint = "https://oauth2.googleapis.com/revoke"
)

type GoogleProvider struct {
	config ProviderConfig
}

func NewGoogleProvider(app *config.App) (OAuthProvider, error) {
	clientID := app.Env.GetString("GOOGLE_CLIENT_ID")
	clientSecret := app.Env.GetString("GOOGLE_CLIENT_SECRET")
	redirectURI := app.Env.GetString("GOOGLE_REDIRECT_URI")

	// Validate required fields
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if redirectURI == "" {
		return nil, fmt.Errorf("GOOGLE_REDIRECT_URI is required")
	}

	config := ProviderConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       []string{"email", "profile"},
	}

	return &GoogleProvider{config: config}, nil
}

func (p *GoogleProvider) GetProviderName() string {
	return "google"
}

func (p *GoogleProvider) CreateAuthorizationURL(state string, codeVerifier string) (*url.URL, error) {

	var queryParams url.Values = url.Values{
		"response_type":         {"code"},
		"client_id":             {p.config.ClientID},
		"state":                 {state},
		"code_challenge":        {utils.CreateS256CodeChallenge(codeVerifier)},
		"code_challenge_method": {"S256"},
		"scope":                 {p.scopesString()},
		"access_type":           {"offline"},
		"prompt":                {"consent"},
	}
	if p.config.RedirectURI != "" {
		queryParams.Set("redirect_uri", p.config.RedirectURI)
	}

	authURL, err := url.Parse(googleAuthEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization endpoint: %w", err)
	}

	authURL.RawQuery = queryParams.Encode()
	return authURL, nil
}

func (p *GoogleProvider) ValidateAuthorizationCode(code string, codeVerifier string) (*utils.OAuth2Tokens, error) {

	var queryParams url.Values = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if codeVerifier != "" {
		queryParams.Set("code_verifier", codeVerifier)
	}
	if p.config.RedirectURI != "" {
		queryParams.Set("redirect_uri", p.config.RedirectURI)
	}
	resp_body, err := p.authFetch(googleTokenEndpoint, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to validate authorization code: %w", err)
	}

	tokens, err := utils.NewOAuth2Tokens(resp_body)
	if err != nil {
		return nil, fmt.Errorf("validation response invalid: %w", err)
	}

	return tokens, nil
}

func (p *GoogleProvider) RefreshAccessToken(refreshToken string) (*utils.OAuth2Tokens, error) {
	var queryParams url.Values = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	if p.scopesString() != "" {
		queryParams.Set("scope", p.scopesString())
	}
	resp_body, err := p.authFetch(googleTokenEndpoint, queryParams)
	fmt.Println("refresh token response:", string(resp_body))

	if err != nil {
		return nil, fmt.Errorf("failed to refresh access token: %w", err)
	}

	tokens, err := utils.NewOAuth2Tokens(resp_body)
	if err != nil {
		return nil, fmt.Errorf("refresh token response invalid: %w", err)
	}

	return tokens, nil
}

func (p *GoogleProvider) RevokeToken(token string) error {

	_, err := p.authFetch(googleRevokeEndpoint,
		url.Values{
			"token": {token},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil

}

func (p *GoogleProvider) GetUserInfo(tokens *utils.OAuth2Tokens) (UserInfo, error) {
	tokenResult, err := tokens.GetTokenResult()
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to get token result: %w", err)
	}

	claims, err := utils.DecodeJwt(tokenResult.IDToken)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to decode ID token: %w", err)
	}

	return UserInfo{
		ID:            claims["sub"].(string),
		Email:         claims["email"].(string),
		Name:          claims["name"].(string),
		Picture:       claims["picture"].(string),
		EmailVerified: claims["email_verified"].(bool),
	}, nil
}

func (p *GoogleProvider) authFetch(endpoint string, queryParams url.Values) ([]byte, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	endpointURL.RawQuery = queryParams.Encode()
	if p.config.ClientSecret == "" {
		queryParams.Set("client_id", p.config.ClientID)
	}

	body := strings.NewReader(queryParams.Encode())

	req, err := http.NewRequest("POST", endpointURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "miles-creative")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Size()))

	if p.config.ClientSecret != "" {
		encodedCredentials := utils.EncodeBasicCredentials(p.config.ClientID, p.config.ClientSecret)
		req.Header.Set("Authorization", "Basic "+encodedCredentials)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %d", resp.StatusCode)
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp_body, nil
}

func (p *GoogleProvider) scopesString() string {
	return strings.Join(p.config.Scopes, " ")
}
