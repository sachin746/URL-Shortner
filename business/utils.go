package business

import (
	"net/http"
	"time"

	"URL-Shortner/constants"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GenerateToken(user *models.User) (string, error) {
	secret := viper.GetString("jwt.userLoginsecret")
	if secret == constants.Empty {
		log.Sugar.Error("User login secret is not set in the configuration")
		return "", constants.ErrSomethingWentWrong
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token valid for 72 hours
		"uat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Sugar.Errorf("Failed to sign token: %v", err)
		return "", constants.ErrSomethingWentWrong
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	secret := viper.GetString("jwt.userLoginsecret")
	if secret == constants.Empty {
		log.Sugar.Error("User login secret is not set in the configuration")
		return nil, constants.ErrSomethingWentWrong
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Sugar.Errorf("Failed to parse token: %v", err)
		return nil, constants.ErrInvalidCredentials
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, constants.ErrInvalidCredentials
}

func GetUserIdFromContext(c *gin.Context) (uuid.UUID, bool) {
	username, exists := c.Get("username")
	if !exists || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return uuid.Nil, false
	}
	userID, err := getUserIDByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid hardcoded UUID"})
		return uuid.Nil, false
	}
	return userID, true
}

func getUserIDByUsername(username string) (uuid.UUID, error) {
	var idStr string
	if err := database.Get().Table("users").
		Select("id").
		Where("username = ?", username).
		Scan(&idStr).Error; err != nil {
		log.Sugar.Errorf("Failed to get user ID by username: %v", err)
		return uuid.Nil, err
	}

	return uuid.Parse(idStr)
}
