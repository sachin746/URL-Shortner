package v1

import (
	"net/http"

	"URL-Shortner/business"
	"URL-Shortner/constants"
	"URL-Shortner/log"
	"URL-Shortner/models"

	"github.com/gin-gonic/gin"
)

func HandleShortenURL(c *gin.Context) {
	var urlRequest models.URL

	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		log.Sugar.Errorf("Failed to bind shorten request: %v", err)
		c.JSON(http.StatusBadRequest, constants.ErrBindJSONFailed.SetErr(err))
		return
	}
	// Validate the login request
	if err := urlRequest.Validate(); err != nil {
		log.Sugar.Errorf("Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, constants.ErrOriginalURLisMandatory.SetErr(err))
		return
	}

	UrlResponse, err := business.GetShortenUrl(urlRequest)
	if err != nil {
		log.Sugar.Errorf("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, constants.ErrInvalidCredentials)
		return
	}
	c.JSON(http.StatusOK, UrlResponse)
	log.Sugar.Infof("User %s logged in successfully", UrlResponse.ShortCode)
}
