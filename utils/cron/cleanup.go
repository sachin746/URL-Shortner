package cron

import (
	"context"
	"log"
	"time"

	"URL-Shortner/entities"
	"URL-Shortner/utils/database"
)

func StartCleanupCron(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cleanupOldRecords()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func cleanupOldRecords() {
	db := database.Get()

	deletes, err := db.Where("valid_till < ?", time.Now()).Delete(&entities.ShortenUrl{}).RowsAffected, db.Error
	if err != nil {
		log.Printf("Failed to clean up old records: %v", err)
	} else {
		log.Printf("Cleaned up %d old records", deletes)
	}
}
