package models

import "URL-Shortner/constants"

type URL struct {
	ShortCode        string `json:"short_code"`
	OriginalURL      string `json:"original_url"`
	UserID           string `json:"user_id"`
	ValidForInMonths int64  `json:"valid_for_in_months"`
}

type CustomURL struct {
	ShortCode        string `json:"short_code"`
	OriginalURL      string `json:"original_url"`
	UserID           string `json:"user_id"`
	ValidForInMonths int64  `json:"valid_for_in_months"`
}

func (u *URL) Validate() error {
	if u.OriginalURL == constants.Empty {
		return constants.ErrInvalidOriginalURL
	} else if u.ValidForInMonths < 1 {
		return constants.ErrInvalidValidityPeriod
	}
	return nil
}

func (u *CustomURL) Validate() error {
	if u.OriginalURL == constants.Empty {
		return constants.ErrInvalidOriginalURL
	} else if u.ValidForInMonths < 1 {
		return constants.ErrInvalidValidityPeriod
	} else if u.ShortCode == constants.Empty || len(u.ShortCode) < 4 || len(u.ShortCode) > 10 {
		return constants.ErrInvalidCustomShortCode
	}
	// shortcode can only contain alphanumeric characters
	keyRunes := []rune(u.ShortCode)
	for _, r := range keyRunes {
		if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
			return constants.ErrInvalidCustomShortCode
		}
	}
	return nil
}
