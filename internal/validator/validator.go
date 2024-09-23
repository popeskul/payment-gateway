package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Struct(s interface{}) error {
	return validate.Struct(s)
}

func Field(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

func ValidateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	validCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"GBP": true,
		// other
	}
	return validCurrencies[currency]
}

func init() {
	err := validate.RegisterValidation("currency", ValidateCurrency)
	if err != nil {
		return
	}
}
