package models

import (
	"time"

	"URL-Shortner/constants"
	"URL-Shortner/utils/errors"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"-"`
	Name       string    `json:"name"`
	Email      string    `json:"email" validate:"required,email"`
	Mobile     string    `json:"mobile"`
	Username   string    `json:"username"`
	DOB        string    `json:"dob"`
	ProfileURL string    `json:"profile_image_url"`
	Password   string    `json:"password,omitempty"`
}

func (u *User) Validate() *errors.Error {
	if u.Name == constants.Empty {
		return constants.ErrNameisMandatory
	}
	if u.Email == constants.Empty || !isValidEmail(u.Email) {
		return constants.ErrEmailisMandatory
	}
	if u.Mobile == constants.Empty || !u.isValidMobile() {
		return constants.ErrMobileisMandatory
	}
	if u.Username == constants.Empty {
		return constants.ErrUsernameisMandatory
	}
	dob, err := time.Parse(constants.DobYYYYMMDDDateFormat, u.DOB)
	if err != nil || dob.IsZero() || dob.After(time.Now()) || u.DOB == constants.Empty {
		return constants.ErrDobisMandatory
	}
	if u.Password == constants.Empty {
		return constants.ErrPasswordisMandatory
	}
	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation logic
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' && (i == 0 || i == len(email)-1) {
			return false
		}
	}
	return at == 1
}

func (u *User) isValidMobile() bool {
	// Simple mobile validation logic
	if len(u.Mobile) < 10 || len(u.Mobile) > 15 {
		return false
	}
	for _, char := range u.Mobile {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
