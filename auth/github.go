package auth

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var GithubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

func GetGithubOAuthState() string {
	return viper.GetString("github.oauth_secret")
}

func GetGithubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  viper.GetString("github.redirect_url"),
		ClientID:     viper.GetString("github.client_id"),
		ClientSecret: viper.GetString("github.client_secret"),
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: GithubEndpoint,
	}
}
