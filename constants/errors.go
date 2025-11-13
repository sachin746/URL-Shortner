package constants

import "habit-tracker/utils/errors"

var (
	ErrUserAlreadyExists   = &errors.Error{Code: "USER-ALREADY-EXISTS", Status: 400, Message: "User already exists", Err: nil}
	ErrEmailAlreadyExists  = &errors.Error{Code: "EMAIL-ALREADY-EXISTS", Status: 400, Message: "Email already exists", Err: nil}
	ErrMobileAlreadyExists = &errors.Error{Code: "MOBILE-ALREADY-EXISTS", Status: 400, Message: "Mobile already exists", Err: nil}
	ErrDBError             = &errors.Error{Code: "DB-ERROR", Status: 500, Message: "Failed to create user", Err: nil}
	ErrNameisMandatory     = &errors.Error{Code: "NAME-IS-MANDATORY", Status: 400, Message: "Name is mandatory", Err: nil}
	ErrEmailisMandatory    = &errors.Error{Code: "EMAIL-IS-MANDATORY", Status: 400, Message: "Email is mandatory", Err: nil}
	ErrPasswordisMandatory = &errors.Error{Code: "PASSWORD-IS-MANDATORY", Status: 400, Message: "Password is mandatory", Err: nil}
	ErrMobileisMandatory   = &errors.Error{Code: "MOBILE-IS-MANDATORY", Status: 400, Message: "Mobile is mandatory", Err: nil}
	ErrDobisMandatory      = &errors.Error{Code: "DOB-IS-MANDATORY", Status: 400, Message: "Date of Birth is invalid", Err: nil}
	ErrUsernameisMandatory = &errors.Error{Code: "USERNAME-IS-MANDATORY", Status: 400, Message: "Username is mandatory", Err: nil}

	ErrBindJSONFailed     = &errors.Error{Code: "BIND-JSON-FAILED", Status: 400, Message: "Failed to bind JSON", Err: nil}
	ErrSomethingWentWrong = &errors.Error{Code: "SOMETHING-WENT-WRONG", Status: 500, Message: "Something went wrong", Err: nil}

	ErrUserNotFound     = &errors.Error{Code: "USER-NOT-FOUND", Status: 404, Message: "User not found", Err: nil}
	ErrValidationFailed = &errors.Error{Code: "VALIDATION-FAILED", Status: 400, Message: "Validation failed", Err: nil}

	ErrInvalidCredentials         = &errors.Error{Code: "INVALID-CREDENTIALS", Status: 401, Message: "Invalid Credentials", Err: nil}
	ErrInternalServerError        = &errors.Error{Code: "INTERNAL-SERVER-ERROR", Status: 500, Message: "Internal server error", Err: nil}
	ErrUsernameOrEmailisMandatory = &errors.Error{Code: "USERNAME-OR-EMAIL-IS-MANDATORY", Status: 400, Message: "Either username or email is mandatory", Err: nil}
	ErrUnauthorized               = &errors.Error{Code: "UNAUTHORIZED", Status: 401, Message: "Unauthorized", Err: nil}

	ErrOwnerRequired	   = &errors.Error{Code: "OWNER-REQUIRED", Status: 403, Message: "Only owner can perform this action", Err: nil}
	ErrOwnerCanAccept = &errors.Error{Code: "OWNER-CAN-ACCEPT", Status: 403, Message: "Only owner can accept other users", Err: nil}
	ErrHabitNotPublic = &errors.Error{Code: "HABIT-NOT-PUBLIC", Status: 403, Message: "Habit is not public, members can only be added by invitation", Err: nil}
)
