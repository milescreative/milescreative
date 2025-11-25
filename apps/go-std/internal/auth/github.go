package auth

import (
	"encoding/json"
	"fmt"
	"go-std/internal/config"
	"go-std/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	githubAuthEndpoint  = "https://github.com/login/oauth/authorize"
	githubTokenEndpoint = "https://github.com/login/oauth/access_token"
	githubUserEndpoint  = "https://api.github.com/user"
	githubEmailEndpoint = "https://api.github.com/user/emails"
)

type GitHubProvider struct {
	config ProviderConfig
}

func NewGitHubProvider(app *config.App) (OAuthProvider, error) {
	clientID := app.Env.GetString("GITHUB_CLIENT_ID")
	clientSecret := app.Env.GetString("GITHUB_CLIENT_SECRET")
	redirectURI := app.Env.GetString("GITHUB_REDIRECT_URI")

	// Validate required fields
	if clientID == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID is required")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_SECRET is required")
	}
	if redirectURI == "" {
		return nil, fmt.Errorf("GITHUB_REDIRECT_URI is required")
	}

	config := ProviderConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       []string{"user:email", "read:user"},
	}

	return &GitHubProvider{config: config}, nil
}

func (p *GitHubProvider) GetProviderName() string {
	return "github"
}

func (p *GitHubProvider) CreateAuthorizationURL(state string, codeVerifier string) (*url.URL, error) {
	queryParams := url.Values{
		"client_id":     {p.config.ClientID},
		"redirect_uri":  {p.config.RedirectURI},
		"scope":         {p.scopesString()},
		"state":         {state},
		"response_type": {"code"},
	}

	authURL, err := url.Parse(githubAuthEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization endpoint: %w", err)
	}

	authURL.RawQuery = queryParams.Encode()
	return authURL, nil
}

func (p *GitHubProvider) ValidateAuthorizationCode(code string, codeVerifier string) (*utils.OAuth2Tokens, error) {
	data := url.Values{
		"client_id":     {p.config.ClientID},
		"client_secret": {p.config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {p.config.RedirectURI},
	}

	respBody, err := p.makeTokenRequest(githubTokenEndpoint, data)
	if err != nil {
		return nil, fmt.Errorf("failed to validate authorization code: %w", err)
	}

	tokens, err := utils.NewOAuth2Tokens(respBody)
	if err != nil {
		return nil, fmt.Errorf("validation response invalid: %w", err)
	}

	return tokens, nil
}

func (p *GitHubProvider) RefreshAccessToken(refreshToken string) (*utils.OAuth2Tokens, error) {
	// GitHub doesn't support refresh tokens, access tokens don't expire
	return nil, fmt.Errorf("github does not support token refresh")
}

func (p *GitHubProvider) RevokeToken(token string) error {
	// GitHub token revocation requires differxent approach
	return fmt.Errorf("github token revocation not implemented")
}

func (p *GitHubProvider) GetUserInfo(tokens *utils.OAuth2Tokens) (UserInfo, error) {
	accessToken, err := tokens.AccessToken()
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to get access token: %w", err)
	}

	// Get user profile
	userResp, err := p.makeAPIRequest(githubUserEndpoint, accessToken)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(userResp, &githubUser); err != nil {
		return UserInfo{}, fmt.Errorf("failed to parse user response: %w", err)
	}

	userInfo := UserInfo{
		ID:            fmt.Sprintf("%d", githubUser.ID),
		Name:          githubUser.Name,
		Picture:       githubUser.AvatarURL,
		EmailVerified: false, // GitHub doesn't provide email verification status
	}

	// If email is not in profile, get from emails endpoint
	if githubUser.Email != "" {
		userInfo.Email = githubUser.Email
		userInfo.EmailVerified = true
	} else {
		email, verified, err := p.getPrimaryEmail(accessToken)
		if err == nil {
			userInfo.Email = email
			userInfo.EmailVerified = verified
		}
	}

	return userInfo, nil
}

func (p *GitHubProvider) getPrimaryEmail(accessToken string) (string, bool, error) {
	emailResp, err := p.makeAPIRequest(githubEmailEndpoint, accessToken)
	if err != nil {
		return "", false, err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(emailResp, &emails); err != nil {
		return "", false, err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, email.Verified, nil
		}
	}

	return "", false, fmt.Errorf("no primary email found")
}

func (p *GitHubProvider) makeTokenRequest(endpoint string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "miles-creative")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *GitHubProvider) makeAPIRequest(endpoint string, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "miles-creative")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API request failed: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *GitHubProvider) scopesString() string {
	return strings.Join(p.config.Scopes, " ")
}
