package infrastructure

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	if err := v.RegisterValidation("allowed_currencies", currencyValidator); err != nil {
		panic(err)
	}

	return &Validator{validator: v}
}

func ValidationErrors(v *Validator, req interface{}) (string, error) {
	if err := v.validator.Struct(req); err != nil {
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
	if err != nil {
		log.Println("error validating req struct in handler", err)
		return errMessage, err
	}

	return "", nil
}

var currencyISOCodes = map[string]struct{}{
	"SGD": {}, "USD": {}, "AUD": {},
}

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	_, ok := currencyISOCodes[currency]
	return ok
}
