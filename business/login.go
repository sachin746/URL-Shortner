package business

import (
	"habit-tracker/constants"
	"habit-tracker/entities"
	"habit-tracker/log"
	"habit-tracker/models"
	"habit-tracker/utils/database"
	"habit-tracker/utils/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginUser(login models.Login) (*models.LoginResponse, *errors.Error) {
	// Check if the user exists
	user, err := getUserByUsernameOrEmail(database.Get(), login.Username, login.Email)
	if err != nil {
		log.Sugar.Errorf("Failed to get user: %v", err)
		return nil, constants.ErrInvalidCredentials
	}
	// get hashed password from the auth model by username
	hashedPassword, err := getHashedPassword(database.Get(), user.ID)
	if err != nil {
		log.Sugar.Errorf("Failed to get hashed password: %v", err)
		return nil, constants.ErrSomethingWentWrong
	}

	// Check if the password is correct
	if !checkPasswordHash(login.Password, *hashedPassword) {
		log.Sugar.Error("Invalid password")
		return nil, constants.ErrInvalidCredentials
	}

	// Generate a JWT token or session ID (not implemented here)
	token, err := GenerateToken(user)
	if err != nil {
		log.Sugar.Errorf("Failed to generate token: %v", err)
		return nil, constants.ErrInternalServerError
	}

	return &models.LoginResponse{
		Message: "Login successful",
		User:    user,
		Token:   token,
	}, nil
}

func getUserByUsernameOrEmail(db *gorm.DB, username, email string) (*models.User, error) {
	var user models.User
	if username == constants.Empty {
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getHashedPassword(db *gorm.DB, id uuid.UUID) (*string, error) {
	var auth entities.Auth
	if err := db.Where("id = ?", id).First(&auth).Error; err != nil {
		return nil, err
	}
	return auth.PasswordHash, nil
}
