package models

import "habit-tracker/constants"

type GetUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (r *GetUserRequest) Validate() error {
	if r.Username == "" {
		return constants.ErrUsernameisMandatory
	}
	if r.Email == "" {
		return constants.ErrEmailisMandatory
	}
	return nil
}
