package business

import (
	"encoding/json"
	"habit-tracker/auth"
	"habit-tracker/constants"
	"habit-tracker/entities"
	"habit-tracker/log"
	"habit-tracker/models"
	"habit-tracker/utils/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleGithubAuthCallback(ctx *gin.Context) (*models.User, string, error) {
	githubState := ctx.Query("state")
	configSecret := auth.GetGithubOAuthState()
	if githubState != configSecret {
		log.Sugar.Error("Invalid OAuth state")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}

	code := ctx.Query("code")
	if code == "" {
		log.Sugar.Error("Authorization code is missing in the context")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}

	token, err := auth.GetGithubOAuthConfig().Exchange(ctx, code)
	if err != nil {
		log.Sugar.Errorf("Failed to exchange code for token: %v", err)
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(err)
	}

	client := auth.GetGithubOAuthConfig().Client(ctx, token)
	userReq, err := GetGithubUserInfo(client)
	if err != nil {
		log.Sugar.Errorf("Failed to get user info: %v", err)
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(err)
	}
	// Extract user information from the response
	if userReq["email"] == nil {
		log.Sugar.Error("Email not found in Google user info")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}
	userEmail, ok := userReq["email"].(string)
	if !ok || userEmail == "" {
		log.Sugar.Error("Invalid or missing email in Google user info")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}

	// check if user exists in the database
	var user models.User
	existingUser := GetifUserExistsByEmail(database.Get(), userEmail)
	if existingUser != nil {
		log.Sugar.Infof("User already exists: %s", existingUser.Username)
		// Update the existing user with new information
		existingUser.ProfileImageURL = userReq["profile_image_url"].(string)
		if err = database.Get().Save(existingUser).Error; err != nil {
			log.Sugar.Errorf("Failed to update existing user: %v", err)
			return nil, constants.Empty, constants.ErrSomethingWentWrong.SetErr(err)
		}
		log.Sugar.Infof("Updated existing user: %s", existingUser.Username)
		user = models.User{
			ID:         existingUser.ID,
			Name:       existingUser.Name,
			Email:      existingUser.Email,
			ProfileURL: existingUser.ProfileImageURL,
			Username:   existingUser.Username,
			Mobile:     existingUser.Mobile,
			DOB:        existingUser.DOB,
		}
	} else {
		log.Sugar.Info("Creating new user from Google OAuth")
		dbUser := entities.User{
			Name:            userReq["name"].(string),
			Email:           userEmail,
			ProfileImageURL: userReq["profile_image_url"].(string),
		}

		if err = database.Get().Create(&dbUser).Error; err != nil {
			log.Sugar.Errorf("Failed to create new user: %v", err)
			return nil, constants.Empty, constants.ErrSomethingWentWrong.SetErr(err)
		}
		log.Sugar.Infof("Created new user with id: %s", user.ID)
		user = models.User{
			ID:         dbUser.ID,
			Name:       dbUser.Name,
			Email:      dbUser.Email,
			ProfileURL: dbUser.ProfileImageURL,
		}
	}
	githubID, ok := userReq["id"].(float64)
	if !ok || githubID == 0 {
		log.Sugar.Error("Invalid or missing Github ID in user info")
	}
	idStr := strconv.FormatInt(int64(githubID), 10)

	// Create or update the auth record
	authRecord := entities.Auth{
		ID:    user.ID,
		Email: &user.Email,
	}
	if err = database.Get().Where("id = ?", user.ID).Assign(entities.Auth{
		GithubID: &idStr,
	}).FirstOrCreate(&authRecord).Error; err != nil {
		log.Sugar.Errorf("Failed to create or update auth record: %v", err)
		return nil, constants.Empty, constants.ErrSomethingWentWrong.SetErr(err)
	}
	log.Sugar.Infof("Auth record created or updated for user: %s", user.Email)

	tokenString, err := GenerateToken(&user)
	if err != nil {
		log.Sugar.Errorf("Failed to generate token: %v", err)
		return nil, constants.Empty, constants.ErrSomethingWentWrong.SetErr(err)
	}

	return &user, tokenString, nil
}

func GetGithubUserInfo(ctx *http.Client) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Sugar.Errorf("Failed to create request: %v", err)
		return nil, constants.ErrSomethingWentWrong.SetErr(err)
	}

	resp, err := ctx.Do(req)
	if err != nil {
		log.Sugar.Errorf("Failed to get user info: %v", err)
		return nil, constants.ErrSomethingWentWrong.SetErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Sugar.Errorf("Failed to get user info: %s", resp.Status)
		return nil, constants.ErrSomethingWentWrong.SetErr(nil)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Sugar.Errorf("Failed to decode user info: %v", err)
		return nil, constants.ErrSomethingWentWrong.SetErr(err)
	}

	var userEmail string
	if email, ok := userInfo["email"].(string); ok && email != "" {
		userEmail = email
	} else {
		// fallback: fetch email from /user/emails
		req2, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		if err != nil {
			log.Sugar.Errorf("Failed to create email request: %v", err)
			return nil, constants.ErrSomethingWentWrong.SetErr(err)
		}
		req2.Header.Set("Accept", "application/vnd.github+json")
		req2.Header.Set("Authorization", req.Header.Get("Authorization")) // reuse token

		resp2, err := ctx.Do(req2)
		if err != nil {
			log.Sugar.Errorf("Failed to get emails from GitHub: %v", err)
			return nil, constants.ErrSomethingWentWrong.SetErr(err)
		}
		defer resp2.Body.Close()

		var emails []map[string]interface{}
		if err := json.NewDecoder(resp2.Body).Decode(&emails); err != nil {
			log.Sugar.Errorf("Failed to decode email list: %v", err)
			return nil, constants.ErrSomethingWentWrong.SetErr(err)
		}

		for _, emailEntry := range emails {
			if emailEntry["primary"].(bool) && emailEntry["verified"].(bool) {
				userEmail = emailEntry["email"].(string)
				break
			}
		}

		if userEmail == "" {
			log.Sugar.Error("No verified primary email found")
			return nil, constants.ErrInvalidCredentials.SetErr(nil)
		}
	}

	if userInfo["id"] == nil {
		log.Sugar.Error("ID not found in Github user info")
		return nil, constants.ErrInvalidCredentials.SetErr(nil)
	}

	userInfo["email"] = userEmail
	userInfo["profile_image_url"] = userInfo["avatar_url"]
	return userInfo, nil
}
