package validation

import "github.com/go-playground/validator/v10"

type ValidatorConfigFunc func(*validator.Validate) error
