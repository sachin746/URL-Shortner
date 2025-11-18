package entities

import "time"

type ShortenUrl struct {
	Id          uint      `gorm:"column:id;primaryKey;"`
	ShortCode   string    `gorm:"column:short_code;type:varchar(10);primaryKey"`
	OriginalURL string    `gorm:"column:original_url;type:text;not null"`
	UserID      string    `gorm:"column:user_id;"`
	Created_at  time.Time `gorm:"column:created_at;autoCreateTime:nano;not null"`
	ValidTill   time.Time `gorm:"column:valid_till;not null"`
}

func (ShortenUrl) TableName() string {
	return "shorten_urls"
}
