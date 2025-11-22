package entities

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null"`            // id SMALLINT PRIMARY KEY
	PasswordHash *string   `gorm:"type:text"`                                // nullable TEXT
	GoogleID     *string   `gorm:"column:google_id;type:varchar(100)"`       // nullable, unique
	Email        *string   `gorm:"type:varchar(100)"`                        // nullable
	GithubID     *string   `gorm:"column:github_id;type:varchar"`            // nullable
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // timestamp default current
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // timestamp default current
}

func (Auth) TableName() string {
	return "auth"
}
