package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func SuccessResponse(data any, msg string) map[string]any {
	return map[string]any{
		"status":  "success",
		"message": msg,
		"data":    data,
	}
}

func ErrorResponse(err error) map[string]any {
	if err == nil {
		return nil
	}

	return map[string]any{
		"status":  "error",
		"message": err.Error(),
	}
}

// convert struct to validation error and then return first error message
func ValidationErrorResponse(err error) map[string]any {
	if err == nil {
		return nil
	}

	// Convert the error to a validation error if possible
	if ve, ok := err.(validator.ValidationErrors); ok {
		return map[string]any{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": ve[0].Translate(nil), // Use a translator if needed
		}
	}
	return map[string]any{
		"status":  "error",
		"code":    http.StatusBadRequest,
		"message": err.Error(),
	}
}

// TODO: Add more fields to the validation error response
type validationError struct {
	Namespace       string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}
