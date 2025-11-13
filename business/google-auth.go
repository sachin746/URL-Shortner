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

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func GetGoogleGithubAuthURL(ctx *gin.Context, app string) (string, error) {
	var url string
	if app == constants.Google {
		secret := auth.GetGoogleOAuthState()
		url = auth.GetGoogleOAuthConfig().AuthCodeURL(secret, oauth2.AccessTypeOffline)
		if url == constants.Empty {
			log.Sugar.Error("Failed to generate Google OAuth URL")
			return "", constants.ErrSomethingWentWrong.SetErr(nil)
		}
	} else if app == constants.Github {
		secret := auth.GetGithubOAuthState()
		url = auth.GetGithubOAuthConfig().AuthCodeURL(secret, oauth2.AccessTypeOffline)
		if url == constants.Empty {
			log.Sugar.Error("Failed to generate Github OAuth URL")
			return "", constants.ErrSomethingWentWrong.SetErr(nil)
		}
	} else {
		log.Sugar.Error("Invalid application type provided for OAuth")
		return "", constants.ErrInvalidCredentials.SetErr(nil)
	}
	return url, nil
}

func HandleGoogleAuthCallback(ctx *gin.Context) (*models.User, string, error) {
	code := ctx.Query("code")
	if code == "" {
		log.Sugar.Error("Authorization code is missing in the context")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}

	token, err := auth.GetGoogleOAuthConfig().Exchange(ctx, code)
	if err != nil {
		log.Sugar.Errorf("Failed to exchange code for token: %v", err)
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(err)
	}

	client := auth.GetGoogleOAuthConfig().Client(ctx, token)
	userReq, err := getGoogleUserInfo(client)
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

	googleID, ok := userReq["sub"].(string)
	if !ok || googleID == "" {
		log.Sugar.Error("Invalid or missing Google ID in user info")
		return nil, constants.Empty, constants.ErrInvalidCredentials.SetErr(nil)
	}
	// check if user exists in the database
	var user models.User
	existingUser := GetifUserExistsByEmail(database.Get(), userEmail)
	if existingUser != nil {
		log.Sugar.Infof("User already exists: %s", existingUser.Username)
		// Update the existing user with new information
		existingUser.ProfileImageURL = userReq["picture"].(string)
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
		googleUser := entities.User{
			Name:            userReq["name"].(string),
			Email:           userEmail,
			ProfileImageURL: userReq["picture"].(string),
		}

		if err := database.Get().Create(&googleUser).Error; err != nil {
			log.Sugar.Errorf("Failed to create new user: %v", err)
			return nil, constants.Empty, constants.ErrSomethingWentWrong.SetErr(err)
		}
		log.Sugar.Infof("Created new user with id: %s", user.ID)
		user = models.User{
			ID:         googleUser.ID,
			Name:       googleUser.Name,
			Email:      googleUser.Email,
			ProfileURL: googleUser.ProfileImageURL,
		}
	}
	authRecord := entities.Auth{
		Email: &user.Email,
		ID:    user.ID,
	}
	if err = database.Get().Where("id = ?", user.ID).Assign(entities.Auth{
		GoogleID: &googleID,
		Email:    &user.Email,
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

func getGoogleUserInfo(client *http.Client) (map[string]interface{}, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Sugar.Errorf("Failed to get user info: %v", err)
		return nil, constants.ErrInvalidCredentials.SetErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Sugar.Errorf("Failed to get user info: %s", resp.Status)
		return nil, constants.ErrInvalidCredentials.SetErr(nil)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Sugar.Errorf("Failed to decode user info: %v", err)
		return nil, constants.ErrInvalidCredentials.SetErr(err)
	}
	return userInfo, nil
}

func GetifUserExistsByEmail(db *gorm.DB, email string) *entities.User {
	var user entities.User
	if err := db.Table("users").Model(&entities.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Sugar.Errorf("Failed to check if user exists: %v", err)
			return nil
		}
		log.Sugar.Info("User does not exist")
		return nil
	}
	return &user
}
