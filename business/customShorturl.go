package business

import (
	"context"
	"time"

	"URL-Shortner/constants"
	"URL-Shortner/entities"
	"URL-Shortner/models"
	"URL-Shortner/utils/database"
)

func CreateCustomShortenUrl(c context.Context, urlRequest models.CustomURL) (models.URL, error) {
	urlentity := entities.ShortenUrl{
		ShortCode:   urlRequest.ShortCode,
		OriginalURL: urlRequest.OriginalURL,
		UserID:      urlRequest.UserID,
	}
	err := database.Get().Where("short_code = ?", urlentity.ShortCode).Table(urlentity.TableName()).First(&urlentity).Error
	if err == nil {
		return models.URL{}, constants.ErrCustomShortCodeAlreadyExists
	}

	validTill := time.Now().AddDate(0, int(urlRequest.ValidForInMonths), 0)

	err = database.Get().Create(&entities.ShortenUrl{
		ShortCode:   urlRequest.ShortCode,
		OriginalURL: urlRequest.OriginalURL,
		UserID:      urlRequest.UserID,
		ValidTill:   validTill,
	}).Error
	if err != nil {
		return models.URL{}, err
	}

	return models.URL{
		ShortCode:        urlRequest.ShortCode,
		OriginalURL:      urlRequest.OriginalURL,
		UserID:           urlRequest.UserID,
		ValidForInMonths: urlRequest.ValidForInMonths,
	}, nil
}
