package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"regexp"
)

const (
	usernamePattern = `^[a-zA-Z0-9_]{5,16}$`
)

var usernameRegexp = regexp.MustCompile(usernamePattern)
var hasLetterRegexp = regexp.MustCompile(`[a-zA-Z]`)

func ProvideValidatorConfigFunc() validation.ValidatorConfigFunc {
	return func(vld *validator.Validate) error {
		return vld.RegisterValidation("username", usernameValidation)
	}
}

func usernameValidation(field validator.FieldLevel) bool {
	username := field.Field().String()
	return usernameRegexp.MatchString(username) &&
		hasLetterRegexp.MatchString(username)
}
