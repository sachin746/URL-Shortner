package entities

import (
	"time"

	"github.com/google/uuid"
)

type Click struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UrlID     uint      `gorm:"column:url_id;not null;index"`
	ClickedAt time.Time `gorm:"column:clicked_at;autoCreateTime"`
	IPAddress string    `gorm:"column:ip_address;type:varchar(45)"`
	Country   string    `gorm:"column:country;type:varchar(100)"`
}

func (Click) TableName() string {
	return "clicks"
}
