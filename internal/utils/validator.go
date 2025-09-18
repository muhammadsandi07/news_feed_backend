package utils

import (
	"news-feed/internal/middleware"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errors[field] = "is required"
		case "min":
			errors[field] = "must be at least " + e.Param() + " characters"
		case "max":
			errors[field] = "must be at most " + e.Param() + " characters"
		default:
			errors[field] = "is invalid"
		}
	}
	return middleware.UnprocessableEntity("validation failed").WithDetails(errors)
}
