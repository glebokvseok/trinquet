package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
)

func ProvideValidatorConfigFunc() validation.ValidatorConfigFunc {
	return func(vld *validator.Validate) error {
		return nil
	}
}
