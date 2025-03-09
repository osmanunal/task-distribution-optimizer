package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data interface{}) []string {
	var errors []string
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, formatValidationError(err))
		}
	}
	return errors
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " alanı zorunludur"
	case "min":
		return err.Field() + " alanı minimum " + err.Param() + " karakter olmalıdır"
	case "max":
		return err.Field() + " alanı maksimum " + err.Param() + " karakter olmalıdır"
	default:
		return err.Field() + " alanı geçersiz"
	}
}

func ValidateLatitude(fl validator.FieldLevel) bool {
	latitude := fl.Field().Float()
	return latitude >= -90 && latitude <= 90
}

func ValidateLongitude(fl validator.FieldLevel) bool {
	longitude := fl.Field().Float()
	return longitude >= -180 && longitude <= 180
}

func init() {
}
