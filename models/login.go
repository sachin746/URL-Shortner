package models

import (
	"URL-Shortner/constants"
)

type Login struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response structure for a successful login
type LoginResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
	Token   string `json:"token"` // JWT token or session ID
}

// Validate checks if the login request is valid
func (l *Login) Validate() error {
	if l.Username == "" && l.Email == "" {
		return constants.ErrUsernameOrEmailisMandatory.SetErr(nil)
	}
	return nil
}
