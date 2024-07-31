package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func CustomErrorMessages(err error) interface{} {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return "invalid input"
	}

	errors := make([]string, 0)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errors = append(errors, fmt.Sprintf("%s is required", field))
		case "min":
			errors = append(errors, fmt.Sprintf("%s must be at least %s characters long", field, err.Param()))
		case "max":
			errors = append(errors, fmt.Sprintf("%s must be at most %s characters long", field, err.Param()))
		case "email":
			errors = append(errors, "invalid email")
		case "phone":
			errors = append(errors, "invalid phone number")
		default:
			errors = append(errors, fmt.Sprintf("%s is invalid", field))
		}
	}

	if len(errors) > 1 {
		return errors
	} else {
		return errors[0]
	}
}
