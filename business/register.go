package business

import (
	"habit-tracker/constants"
	"habit-tracker/entities"
	"habit-tracker/log"
	"habit-tracker/models"
	"habit-tracker/utils/database"
	"habit-tracker/utils/errors"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(user models.User) (*models.LoginResponse, error) {
	if CheckUserExists(database.Get(), user.Username) {
		return nil, constants.ErrUserAlreadyExists
	}

	if checkEmailExists(database.Get(), user.Email) {
		return nil, constants.ErrEmailAlreadyExists
	}

	if checkMobileExists(database.Get(), user.Mobile) {
		return nil, constants.ErrMobileAlreadyExists
	}

	// save the user and hash the password
	userEntity := entities.User{
		ID:       uuid.New(),
		Name:     user.Name,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Username: user.Username,
		DOB:      user.DOB,
	}
	savedUser, token, err := saveUserAndPassword(userEntity, user.Password)
	if err != nil {
		log.Sugar.Errorf("Failed to save user: %v", err)
		return nil, constants.ErrSomethingWentWrong
	}

	return &models.LoginResponse{
		Message: "User registered successfully",
		User:    savedUser,
		Token:   *token,
	}, nil
}

func saveUserAndPassword(user entities.User, pass string) (*models.User, *string, *errors.Error) {
	// begin transaction
	tx := database.Get().Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, nil, &errors.Error{
			Code:    constants.ErrDBError.Code,
			Status:  constants.ErrDBError.Status,
			Message: constants.ErrDBError.Message,
			Err:     err,
		}
	}
	savedUser := models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Username: user.Username,
		DOB:      user.DOB,
	}

	// hash the password
	hashedPassword, err := hashPassword(pass)
	if err != nil {
		tx.Rollback()
		log.Sugar.Errorf("Failed to hash password: %v", err)
		return nil, nil, constants.ErrSomethingWentWrong
	}

	auth := entities.Auth{
		ID:           savedUser.ID,
		Email:        &savedUser.Email,
		PasswordHash: &hashedPassword,
	}
	if err := tx.Create(&auth).Error; err != nil {
		tx.Rollback()
		log.Sugar.Errorf("Failed to save auth: %v", err)
		return nil, nil, constants.ErrSomethingWentWrong
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Sugar.Errorf("Failed to commit transaction: %v", err)
		return nil, nil, constants.ErrSomethingWentWrong
	}
	token, err := GenerateToken(&savedUser)
	if err != nil {
		log.Sugar.Errorf("Failed to generate token: %v", err)
		return nil, nil, constants.ErrSomethingWentWrong
	}

	return &savedUser, &token, nil
}

func hashPassword(password string) (string, error) {
	defaultCost := viper.GetInt("bcrypt.cost")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		log.Sugar.Errorf("Failed to hash password: %v", err)
		return constants.Empty, constants.ErrSomethingWentWrong
	}
	return string(bytes), nil
}

func CheckUserExists(db *gorm.DB, username string) bool {
	var count int64
	if err := db.Table("users").Model(&entities.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false // Error occurred, assume user does not exist
	}
	return count > 0
}

func checkEmailExists(db *gorm.DB, email string) bool {
	var count int64
	if err := db.Table("users").Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false // Error occurred, assume email does not exist
	}
	return count > 0
}

func checkMobileExists(db *gorm.DB, mobile string) bool {
	var count int64
	if err := db.Table("users").Model(&entities.User{}).Where("mobile = ?", mobile).Count(&count).Error; err != nil {
		return false // Error occurred, assume mobile does not exist
	}
	return count > 0
}
