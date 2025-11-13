package models

import "URL-Shortner/constants"

type URL struct {
	ShortCode   string
	OriginalURL string
	UserID      uint
	CreatedAt   int64
	ExpiresAt   int64
}

func (u *URL) Validate() error {
	if u.OriginalURL == "" {
		return constants.ErrInvalidOriginalURL
	}
	return nil
}
