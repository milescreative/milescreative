package auth

import (
	"fmt"
	"go-std/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	authorizationEndpoint   = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint           = "https://oauth2.googleapis.com/token"
	tokenRevocationEndpoint = "https://oauth2.googleapis.com/revoke"
)

type GoogleOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scopes       []string
}

func NewGoogleOAuth(clientID string, clientSecret string, redirectURI string, scopes []string) *GoogleOAuth {
	return &GoogleOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		scopes:       scopes,
	}
}

func (o *GoogleOAuth) CreateAuthorizationURLWithPKCE(state string, codeVerifier string) (*url.URL, error) {

	var queryParams url.Values = url.Values{
		"response_type":         {"code"},
		"client_id":             {o.clientID},
		"state":                 {state},
		"code_challenge":        {utils.CreateS256CodeChallenge(codeVerifier)},
		"code_challenge_method": {"S256"},
		"scope":                 {o.Scopes()},
	}
	if o.redirectURI != "" {
		queryParams.Set("redirect_uri", o.redirectURI)
	}

	url_, err := url.Parse(authorizationEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization endpoint: %w", err)
	}

	url_.RawQuery = queryParams.Encode()
	return url_, nil
}

func (o *GoogleOAuth) ValidateAuthorizationCode(code string, codeVerifier string) (*utils.OAuth2Tokens, error) {

	var queryParams url.Values = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if codeVerifier != "" {
		queryParams.Set("code_verifier", codeVerifier)
	}
	if o.redirectURI != "" {
		queryParams.Set("redirect_uri", o.redirectURI)
	}
	resp_body, err := o.AuthFetch(tokenEndpoint, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to validate authorization code: %w", err)
	}

	tokens, err := utils.NewOAuth2Tokens(resp_body)
	if err != nil {
		return nil, fmt.Errorf("validation response invalid: %w", err)
	}

	return tokens, nil
}

func (o *GoogleOAuth) RefreshAccessToken(refreshToken string) (*utils.OAuth2Tokens, error) {
	var queryParams url.Values = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	if o.Scopes() != "" {
		queryParams.Set("scope", o.Scopes())
	}
	resp_body, err := o.AuthFetch(tokenEndpoint, queryParams)

	if err != nil {
		return nil, fmt.Errorf("failed to refresh access token: %w", err)
	}

	tokens, err := utils.NewOAuth2Tokens(resp_body)
	if err != nil {
		return nil, fmt.Errorf("refresh token response invalid: %w", err)
	}

	return tokens, nil
}

func (o *GoogleOAuth) RevokeToken(token string) error {

	_, err := o.AuthFetch(tokenRevocationEndpoint,
		url.Values{
			"token": {token},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil

}

func (o *GoogleOAuth) AuthFetch(endpoint string, queryParams url.Values) ([]byte, error) {
	url_, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	url_.RawQuery = queryParams.Encode()
	if o.clientSecret == "" {
		queryParams.Set("client_id", o.clientID)
	}

	body := strings.NewReader(queryParams.Encode())

	req, err := http.NewRequest("POST", url_.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "miles-creative")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Size()))

	if o.clientSecret != "" {
		encodedCredentials := utils.EncodeBasicCredentials(o.clientID, o.clientSecret)
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

func (o *GoogleOAuth) Scopes() string {
	var scopeString string
	for _, scope := range o.scopes {
		scopeString += scope + " "
	}
	return scopeString
}
