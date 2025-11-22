package models

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID           uuid.UUID `gorm:"primaryKey;not null"`                      // id SMALLINT PRIMARY KEY
	PasswordHash *string   `gorm:"type:text"`                                // nullable TEXT
	GoogleID     *string   `gorm:"type:varchar(100);"`                       // nullable, unique
	GoogleEmail  *string   `gorm:"type:varchar(100)"`                        // nullable
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // timestamp default current
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // timestamp default current
}
