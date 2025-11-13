package business

import (
	"URL-Shortner/constants"
	"URL-Shortner/entities"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils/database"
)

func GetUser(username string) (*models.User, error) {
	user, err := getUserByUsername(username)
	if err != nil {
		log.Sugar.Errorf("Failed to get user by username %s: %v", username, err)
		return nil, constants.ErrUserNotFound.SetErr(err)
	}

	return user, nil
}

func getUserByUsername(username string) (*models.User, error) {
	var user entities.User
	err := database.Get().Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Sugar.Errorf("Database error while fetching user: %v", err)
		return nil, err
	}
	userModel := models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		DOB:      user.DOB,
		Mobile:   user.Mobile,
	}
	return &userModel, nil
}
