package v1

import (
	"URL-Shortner/business"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils"

	"github.com/gin-gonic/gin"
)

func HandleCustomShortenURL(c *gin.Context) {
	// Implementation for handling custom short URL creation
	var req models.CustomURL
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Sugar.Errorf("Failed to bind request body: %v", err)
		c.JSON(400, utils.ErrorResponse(err))
		return
	}

	if err := req.Validate(); err != nil {
		log.Sugar.Errorf("Validation error: %v", err)
		c.JSON(400, utils.ErrorResponse(err))
		return
	}

	// Call business logic to create custom short URL
	createdURL, err := business.CreateCustomShortenUrl(c, req)
	if err != nil {
		log.Sugar.Errorf("Failed to create custom short URL: %v", err)
		c.JSON(500, utils.ErrorResponse(err))
		return
	}

	c.JSON(201, createdURL)
}
