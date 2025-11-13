package v1

import (
	"net/http"

	"URL-Shortner/models"

	"URL-Shortner/business"
	"URL-Shortner/constants"
	"URL-Shortner/log"

	"github.com/gin-gonic/gin"
)

func HandleLoginUser(c *gin.Context) {
	// Extract login credentials from the request
	var loginRequest models.Login

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Sugar.Errorf("Failed to bind login request: %v", err)
		c.JSON(http.StatusBadRequest, constants.ErrBindJSONFailed.SetErr(err))
		return
	}
	// Validate the login request
	if err := loginRequest.Validate(); err != nil {
		log.Sugar.Errorf("Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, constants.ErrInvalidCredentials.SetErr(err))
		return
	}

	userData, err := business.LoginUser(loginRequest)
	if err != nil {
		log.Sugar.Errorf("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, constants.ErrInvalidCredentials)
		return
	}
	c.JSON(http.StatusOK, userData)
	log.Sugar.Infof("User %s logged in successfully", userData.User)
}
