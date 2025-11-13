package utils

import (

	// "unicode"

	"unicode"

	"github.com/go-playground/validator/v10"
)

// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	// err := Validate.RegisterValidation("validate_empty_string", CustomEmptyStringValidator)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}

// CustomEmptyStringValidator checks if a string is empty or contains only whitespace
func CustomEmptyStringValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return !IsBlank(&value)
}

func IsBlank(strPtr *string) bool {
	if strPtr == nil {
		return true
	}

	str := *strPtr
	strLen := len(str)
	if strLen == 0 {
		return true
	}

	for i := 0; i < strLen; i++ {
		if !unicode.IsSpace(rune(str[i])) {
			return false
		}
	}
	return true
}
