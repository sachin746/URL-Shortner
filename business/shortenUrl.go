package business

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"URL-Shortner/entities"
	"URL-Shortner/log"
	"URL-Shortner/models"
	"URL-Shortner/utils/cache"
	"URL-Shortner/utils/database"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func GetShortenUrl(urlRequest models.URL) (models.URL, error) {
	for range 3 {
		id, err := getUrlId(database.Get())
		if err != nil {
			log.Sugar.Errorf("Failed to generate URL ID: %v", err)
			return models.URL{}, err
		}
		println("Generated URL ID:", id)
		shortCode, err := generateShortCode(id)
		if err != nil {
			log.Sugar.Errorf("Failed to generate short code: %v", err)
			return models.URL{}, err
		}
		println("Generated Short Code:", shortCode)
		validTill := time.Now().AddDate(0, int(urlRequest.ValidForInMonths), 0)

		// Save the shortened URL to the database
		shortenUrl := entities.ShortenUrl{
			ShortCode:   shortCode,
			OriginalURL: urlRequest.OriginalURL,
			UserID:      urlRequest.UserID,
			ValidTill:   validTill,
		}

		result := database.Get().Create(&shortenUrl)
		if result.Error == nil {
			urlRequest.ShortCode = shortCode
			return urlRequest, nil
		}
	}
	return models.URL{}, errors.New("failed to generate unique short code after multiple attempts")
}

func getUrlId(db *gorm.DB) (uint, error) {
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
		return 0, errors.New("failed to get URL ID generator")
	}

	currentId := urlIdGenerator.CurrentID

	// increment current id by 1
	urlIdGenerator.CurrentID += 1
	tx.Save(&urlIdGenerator)
	tx.Commit()

	return currentId, nil
}

func generateShortCode(id uint) (string, error) {
	mappings := make(map[string]string)
	mappings = viper.GetStringMapString("characherMapping")
	// Convert the ID to a short code using base62 encoding

	binary := strconv.FormatUint(uint64(id), 2)
	paddedbinaryString := fmt.Sprintf("%048s", binary)
	log.Sugar.Info("Padded Binary String:", paddedbinaryString)
	log.Sugar.Info("Length of Padded Binary String:", len(paddedbinaryString))
	var shortCode string
	for i := 0; i < len(paddedbinaryString); i += 6 {
		chunk := paddedbinaryString[i : i+6]
		// binary string → int
		n, err := strconv.ParseInt(chunk, 2, 64)
		if err != nil {
			log.Sugar.Errorf("Failed to parse binary string: %v", err)
			return "", err
		}
		fmt.Println("Chunk:", chunk, "Int:", n)
		shortCode += mappings[fmt.Sprintf("%d", n)]
		log.Sugar.Info("Current Short Code:", shortCode)
	}
	return shortCode, nil
}

func GetOriginalURL(ctx context.Context, shortCode string) (models.URL, error) {
	OriginalURL, err := cache.GetRedisCache().Get(ctx, shortCode)
	if err != nil && err != redis.Nil {
		log.Sugar.Errorf("Failed to get URL from cache: %v", err)
	} else if err == nil {
		urlResponse := models.URL{
			OriginalURL: OriginalURL,
			ShortCode:   shortCode,
		}
		log.Sugar.Info("Cache hit for short code:", shortCode)
		return urlResponse, nil
	}
	var shortenUrl entities.ShortenUrl
	result := database.Get().Where("short_code = ?", shortCode).First(&shortenUrl)
	if result.Error != nil {
		log.Sugar.Errorf("Failed to retrieve original URL: %v", result.Error)
		return models.URL{}, result.Error
	}

	urlResponse := models.URL{
		OriginalURL:      shortenUrl.OriginalURL,
		ShortCode:        shortenUrl.ShortCode,
		UserID:           shortenUrl.UserID,
		ValidForInMonths: int64(shortenUrl.ValidTill.Month()),
	}

	err = cache.GetRedisCache().Set(ctx, shortCode, shortenUrl.OriginalURL, 1*time.Hour)
	if err != nil {
		log.Sugar.Errorf("Failed to set URL in cache: %v", err)
	}
	return urlResponse, nil
}
