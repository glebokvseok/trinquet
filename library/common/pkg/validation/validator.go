package validation

import "github.com/go-playground/validator/v10"

func provideValidator(
	configureValidator ValidatorConfigFunc,
) (*validator.Validate, error) {
	vld := validator.New()

	err := configureValidator(vld)
	if err != nil {
		return nil, err
	}

	return vld, nil
}
