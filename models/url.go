package models

import "URL-Shortner/constants"

type URL struct {
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
