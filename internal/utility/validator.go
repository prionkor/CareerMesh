package utility

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	if err := validate.RegisterValidation("notblank", func(field validator.FieldLevel) bool {
		return field.Field().Kind() == reflect.String && strings.TrimSpace(field.Field().String()) != ""
	}); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("model_valid", func(field validator.FieldLevel) bool {
		value, ok := field.Field().Interface().(interface{ IsValid() bool })
		return ok && value.IsValid()
	}); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("absolute_http_url", func(field validator.FieldLevel) bool {
		parsed, err := url.Parse(field.Field().String())
		if err != nil || !parsed.IsAbs() || parsed.Hostname() == "" {
			return false
		}
		switch strings.ToLower(parsed.Scheme) {
		case "http", "https":
			return true
		default:
			return false
		}
	}); err != nil {
		panic(err)
	}
}

func ValidateStruct(value any) error {
	return validate.Struct(value)
}
