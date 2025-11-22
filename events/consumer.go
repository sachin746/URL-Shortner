package events

import (
	"URL-Shortner/entities"
	"URL-Shortner/log"
	"URL-Shortner/utils/database"

	"gorm.io/gorm"
)

// StartClickWorker starts the background worker to process click events
func StartClickWorker() {
	go func() {
		for event := range ClickQueue {
			processClickEvent(event)
		}
	}()
}

func processClickEvent(event ClickEvent) {
	db := database.Get()
	if db == nil {
		log.Sugar.Error("Database connection is nil, cannot process click event")
		return
	}

	// Find the URL entity to get its ID
	var urlEntity entities.ShortenUrl
	if err := db.Where("short_code = ?", event.ShortCode).First(&urlEntity).Error; err != nil {
		log.Sugar.Errorf("Failed to find URL for short code %s: %v", event.ShortCode, err)
		return
	}

	// Create Click entry
	click := entities.Click{
		UrlID:     urlEntity.Id,
		ClickedAt: event.ClickedAt,
		IPAddress: event.IPAddress,
		// Country:   "Unknown", // GeoIP lookup can be added here
	}

	// Transaction to ensure consistency
	err := db.Transaction(func(tx *gorm.DB) error {
		// Insert click record
		if err := tx.Create(&click).Error; err != nil {
			return err
		}

		// Increment click count
		if err := tx.Model(&urlEntity).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Sugar.Errorf("Failed to process click event for %s: %v", event.ShortCode, err)
	}
}
