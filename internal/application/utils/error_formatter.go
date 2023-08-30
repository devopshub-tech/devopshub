package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatErrors(err error) []string {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	var errorMessages []string
	for _, err := range validationErrs {
		message := fmt.Sprintf("'%s' field validation failed on '%s' tag.", err.Field(), err.Tag())
		errorMessages = append(errorMessages, message)
	}

	return errorMessages
}
