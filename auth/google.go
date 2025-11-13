package auth

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  viper.GetString("google.redirect_url"),
		ClientID:     viper.GetString("google.client_id"),
		ClientSecret: viper.GetString("google.client_secret"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},

		Endpoint: google.Endpoint,
	}
}

func GetGoogleOAuthState() string {
	return viper.GetString("google.oauth_secret")
}
