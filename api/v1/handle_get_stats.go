package v1

import (
	"net/http"

	"URL-Shortner/entities"
	"URL-Shortner/log"
	"URL-Shortner/utils/database"

	"github.com/gin-gonic/gin"
)

type StatsResponse struct {
	ShortCode    string           `json:"short_code"`
	OriginalURL  string           `json:"original_url"`
	ClickCount   int64            `json:"click_count"`
	RecentClicks []entities.Click `json:"recent_clicks"`
}

func HandleGetStats(c *gin.Context) {
	shortCode := c.Param("shortcode")
	db := database.Get()

	var urlEntity entities.ShortenUrl
	if err := db.Where("short_code = ?", shortCode).First(&urlEntity).Error; err != nil {
		log.Sugar.Errorf("Failed to find URL for stats %s: %v", shortCode, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	var recentClicks []entities.Click
	// Fetch last 50 clicks
	if err := db.Where("url_id = ?", urlEntity.Id).Order("clicked_at desc").Limit(50).Find(&recentClicks).Error; err != nil {
		log.Sugar.Errorf("Failed to fetch clicks for %s: %v", shortCode, err)
		// Don't fail the whole request, just return empty clicks
	}

	response := StatsResponse{
		ShortCode:    urlEntity.ShortCode,
		OriginalURL:  urlEntity.OriginalURL,
		ClickCount:   urlEntity.ClickCount,
		RecentClicks: recentClicks,
	}

	c.JSON(http.StatusOK, response)
}
