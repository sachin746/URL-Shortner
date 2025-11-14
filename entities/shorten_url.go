package entities

type ShortenUrl struct {
	Id          uint   `gorm:"column:id;primaryKey;"`
	ShortCode   string `gorm:"column:short_code;type:varchar(10);primaryKey"`
	OriginalURL string `gorm:"column:original_url;type:text;not null"`
	UserID      uint   `gorm:"column:user_id;"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime:nano;not null"`
	ExpiresAt   int64  `gorm:"column:expires_at;not null"`
}

func (ShortenUrl) TableName() string {
	return "shorten_urls"
}
