package v1

import (
	"errors"
	"net/http"
	"time"

	"URL-Shortner/business"
	"URL-Shortner/constants"
	"URL-Shortner/events"
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
	log.Sugar.Infof("Received Request %v", urlRequest)
	err := urlRequest.Validate()
	if err != nil {
		log.Sugar.Errorf("Validation failed: %v", err)
		c.JSON(http.StatusBadRequest, err)
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

func HandleGetURL(c *gin.Context) {
	shortCode := c.Param("shortcode")
	urlResponse, err := business.GetOriginalURL(c, shortCode)
	if err != nil {
		log.Sugar.Errorf("Failed to get original URL: %v", err)
		c.JSON(http.StatusNotFound, errors.New("shortcode not found"))
		return
	}
	// Publish click event asynchronously
	events.PublishClickEvent(events.ClickEvent{
		ShortCode: shortCode,
		IPAddress: c.ClientIP(),
		ClickedAt: time.Now(),
	})

	c.Redirect(http.StatusMovedPermanently, urlResponse.OriginalURL)
	log.Sugar.Infof("Original URL for shortcode %s retrieved successfully", shortCode)
}
