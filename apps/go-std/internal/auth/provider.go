package auth

import (
	"fmt"
	"go-std/internal/config"
	"go-std/internal/utils"
	"net/url"
)

// OAuthProvider defines the interface that all OAuth providers must implement
type OAuthProvider interface {
	// GetProviderName returns the name of the provider (e.g., "google", "github")
	GetProviderName() string

	// CreateAuthorizationURL creates the OAuth authorization URL with PKCE
	CreateAuthorizationURL(state string, codeVerifier string) (*url.URL, error)

	// ValidateAuthorizationCode exchanges the authorization code for tokens
	ValidateAuthorizationCode(code string, codeVerifier string) (*utils.OAuth2Tokens, error)

	// RefreshAccessToken refreshes an access token using a refresh token
	RefreshAccessToken(refreshToken string) (*utils.OAuth2Tokens, error)

	// RevokeToken revokes a token
	RevokeToken(token string) error

	// GetUserInfo extracts user information from the ID token
	GetUserInfo(tokens *utils.OAuth2Tokens) (UserInfo, error)
}

// UserInfo represents standardized user information from OAuth providers
type UserInfo struct {
	ID            string
	Email         string
	Name          string
	Picture       string
	EmailVerified bool
}

// ProviderConfig holds common configuration for OAuth providers
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type ProviderRegistry struct {
	app       *config.App
	providers map[string]ProviderConstructor
}

type ProviderConstructor func(app *config.App) (OAuthProvider, error)

// NewProviderFactory creates a new provider factory
func NewProviderRegistry(app *config.App) (*ProviderRegistry, error) {
	registry := &ProviderRegistry{
		app:       app,
		providers: make(map[string]ProviderConstructor),
	}

	// Register all available providers
	registry.registerProvider("google", NewGoogleProvider)
	registry.registerProvider("github", NewGitHubProvider)

	if err := registry.validateProviders(); err != nil {
		return nil, fmt.Errorf("provider validation failed: %w", err)
	}

	return registry, nil
}

func (r *ProviderRegistry) registerProvider(name string, constructor ProviderConstructor) {
	r.providers[name] = constructor
}

// CreateProvider creates an OAuth provider instance
func (r *ProviderRegistry) CreateProvider(name string) (OAuthProvider, error) {
	constructor, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not supported", name)
	}
	return constructor(r.app)
}

func (r *ProviderRegistry) GetSupportedProviders() []string {
	var providers []string
	for name := range r.providers {
		providers = append(providers, name)
	}
	return providers
}

func (r *ProviderRegistry) validateProviders() error {
	for name, constructor := range r.providers {
		_, err := constructor(r.app)
		if err != nil {
			return fmt.Errorf("provider %s: %w", name, err)
		}
	}
	return nil
}
