package auth

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const CallbackURLPrefix = "http://localhost:3000/api/auth/callback/"

type Provider struct {
	name         string
	gothProvider goth.Provider
}

var Providers = []Provider{
	{name: "google",
		gothProvider: google.New(os.Getenv("OAUTH_KEY"),
			os.Getenv("OAUTH_SECRET"),
			CallbackURLPrefix+"google",
			"email",
			"profile",
			"openid",
		)},
	{name: "github",
		gothProvider: github.New(
			os.Getenv("GITHUB_OAUTH_KEY"),
			os.Getenv("GITHUB_OAUTH_SECRET"),
			CallbackURLPrefix+"github",
			"user:email",
			"read:user",
		)},
}

func SetupProviders() {
	gothProviders := []goth.Provider{}
	for _, provider := range Providers {
		gothProviders = append(gothProviders, provider.gothProvider)
	}
	goth.UseProviders(gothProviders...)
}

func GetProviders() []Provider {
	providers := []Provider{}

	providers = append(providers, Providers...)

	return providers
}

func GetCallbackURL(provider string) string {
	return CallbackURLPrefix + provider
}
