package infrastructure

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func ValidationErrors(validate *validator.Validate, req interface{}) (string, error) {
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return strings.Join(validationErrors, " "), err
	}
	return "", nil
}

func ValidateStruct(req interface{}) (string, error) {
	validate := NewValidator()
	errMessage, err := ValidationErrors(validate, req)
	return errMessage, err
}
