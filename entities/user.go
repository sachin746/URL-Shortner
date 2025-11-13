package entities

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID `json:"-" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string    `json:"name" gorm:"type:varchar(100);not null;column:full_name"`
	Email           string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Mobile          string    `json:"mobile" gorm:"type:varchar(15);unique;not null"`
	Username        string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	DOB             string    `json:"dob" gorm:"type:date;not null"`
	ProfileImageURL string    `json:"profile_image_url" gorm:"type:varchar(255)"`
	PasswordHash    string    `json:"-" gorm:"type:varchar(255)"`
	Coins           int       `json:"coins" gorm:"default:0;not null"`
}

func TableName() string {
	return "users"
}
