package business

import (
	"errors"
	"fmt"

	"URL-Shortner/entities"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils/database"

	"gorm.io/gorm"
)

func GetShortenUrl(urlRequest models.URL) (models.URL, error) {
	id, err := getUrlId(database.Get())
	if err != nil {
		log.Sugar.Errorf("Failed to generate URL ID: %v", err)
		return models.URL{}, err
	}
	println("Generated URL ID:", id)
	shortCode := generateShortCode(uint(id[0]))
	println("Generated Short Code:", shortCode)

	// Save the shortened URL to the database
	shortenUrl := entities.ShortenUrl{
		ShortCode:   shortCode,
		OriginalURL: urlRequest.OriginalURL,
		UserID:      urlRequest.UserID,
		ExpiresAt:   urlRequest.ExpiresAt,
	}

	result := database.Get().Create(&shortenUrl)
	if result.Error != nil {
		log.Sugar.Errorf("Failed to save shortened URL: %v", result.Error)
		return models.URL{}, result.Error
	}

	urlRequest.ShortCode = shortCode
	// Return the shortened URL response
	return urlRequest, nil
}

func getUrlId(db *gorm.DB) (string, error) {
	// get current id from url_id table
	var urlIdGenerator entities.URLIDGenerator
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Statement.Raw("SELECT * FROM url_id where current_id < range_end order by random() limit 1 FOR UPDATE").
		Scan(&urlIdGenerator).
		Error
	if err != nil {
		log.Sugar.Errorf("Failed to get URL ID generator: %v", err)
		tx.Rollback()
		return "", errors.New("failed to get URL ID generator")
	}

	currentId := urlIdGenerator.CurrentID

	// increment current id by 1
	urlIdGenerator.CurrentID += 1
	tx.Save(&urlIdGenerator)
	tx.Commit()

	return fmt.Sprintf("%d", currentId), nil
}

func generateShortCode(id uint) string {
	// Convert the ID to a short code using base62 encoding
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var shortCode []byte
	for id > 0 {
		remainder := id % 62
		shortCode = append([]byte{charset[remainder]}, shortCode...)
		id = id / 62
	}
	return string(shortCode)
}

func GetOriginalURL(shortCode string) (models.URL, error) {
	var shortenUrl entities.ShortenUrl
	result := database.Get().Where("short_code = ?", shortCode).First(&shortenUrl)
	if result.Error != nil {
		log.Sugar.Errorf("Failed to retrieve original URL: %v", result.Error)
		return models.URL{}, result.Error
	}

	urlResponse := models.URL{
		OriginalURL: shortenUrl.OriginalURL,
		ShortCode:   shortenUrl.ShortCode,
		UserID:      shortenUrl.UserID,
		ExpiresAt:   shortenUrl.ExpiresAt,
	}

	return urlResponse, nil
}
