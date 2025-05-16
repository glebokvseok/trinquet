package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"regexp"
)

const (
	passwordPattern = `^[A-Za-z\d!@#$%^&*]{8,20}$`
)

var passwordRegexp = regexp.MustCompile(passwordPattern)
var hasUppercaseRegexp = regexp.MustCompile(`[A-Z]`)
var hasDigitRegexp = regexp.MustCompile(`\d`)
var hasSpecialRegexp = regexp.MustCompile(`[!@#$%^&*]`)

func ProvideValidatorConfigFunc() validation.ValidatorConfigFunc {
	return func(vld *validator.Validate) error {
		return vld.RegisterValidation("password", passwordValidation)
	}
}

func passwordValidation(field validator.FieldLevel) bool {
	password := field.Field().String()
	return passwordRegexp.MatchString(password) &&
		hasUppercaseRegexp.MatchString(password) &&
		hasDigitRegexp.MatchString(password) &&
		hasSpecialRegexp.MatchString(password)
}
