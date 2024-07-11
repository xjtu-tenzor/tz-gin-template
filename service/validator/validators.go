package validator

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

// phone format should be like: (+xx)xxx xxxx xxxx
func phone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	phoneNoSpaces := strings.ReplaceAll(phone, " ", "")

	match, _ := regexp.MatchString(`^(\+\d{1,2})?(\d{11})$`, phoneNoSpaces)
	return match
}

// eamil format should be like: xxx@.xxx(__VA_ARGS__)
func email(fl validator.FieldLevel) bool {
	em, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)
	return regex.MatchString(em)
}
