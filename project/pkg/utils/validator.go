package utils

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateUser(user models.User) []map[string]string {
	var errors []map[string]string
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, errorMessages(err))
		}
	}
	return errors
}

func errorMessages(err validator.FieldError) map[string]string {
	errTag := err.Tag()
	errField := strings.ToLower(err.Field())
	errMsg := make(map[string]string)

	errMsg["field"] = errField

	switch errTag {
	case "required":
		errMsg["message"] = errField + " is required"
	case "email":
		errMsg["message"] = errField + " should be in an email format"
	case "min", "max":
		errMsg["message"] = errField + " " + errTag + " should be " + err.Param()
	}
	return errMsg
}
