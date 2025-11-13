package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"habit-tracker/business"
	"habit-tracker/constants"
	"habit-tracker/log"
	"habit-tracker/utils"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// HandleGoogleAuth godoc
// // @Summary Initiate Google OAuth2 authentication
// // // @Description Redirects the user to Google OAuth2 authentication page
// // // @Tags auth
// // // @Success 302 {string} string "Redirects to Google OAuth2 authentication page"
// // // @Failure 500 {object} constants.Error "Internal server error"
// // // @Router /v1/auth/google [get]
func HandleGoogleAuth(c *gin.Context) {
	// Generate the URL for Google OAuth2 authentication
	url, err := business.GetGoogleGithubAuthURL(c, constants.Google)
	if err != nil {
		log.Sugar.Errorf("Failed to get Google OAuth URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication URL"})
		return
	}
	// Redirect the user to the Google OAuth2 authentication page
	c.Redirect(http.StatusTemporaryRedirect, url)
	log.Sugar.Info("Redirecting to Google OAuth2 authentication page")
}

// HandleGoogleAuthCallback godoc
// // @Summary Handle Google OAuth2 callback
// // // @Description Handles the callback from Google OAuth2 after user authentication
// // // @Tags auth
// // // @Success 200 {object} models.LoginResponse "Login response containing user data and token"
// // // @Failure 401 {object} constants.Error "Unauthorized"
// // // @Failure 500 {object} constants.Error "Internal server error"
// // // @Router /v1/google-auth/callback [get]
func HandleGoogleAuthCallback(c *gin.Context) {
	// Handle the callback from Google OAuth2
	user, token, err := business.HandleGoogleAuthCallback(c)
	if err != nil {
		log.Sugar.Errorf("Failed to handle Google OAuth callback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Google"})
		return
	}

	if user == nil {
		log.Sugar.Error("User not found in Google authentication callback")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(errors.New("user not found")))
		return
	}

	log.Sugar.Infof("User %s authenticated successfully with Google", user.Username)
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Sugar.Infof("Failed to marshal user data: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/oauth/callback?error=user_data_error")
		return
	}

	// Redirect to frontend with token and user data
	frontendURL := fmt.Sprintf("http://localhost:3000/oauth/callback?token=%s&user=%s",
		token,
		url.QueryEscape(string(userJSON)))

	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}
