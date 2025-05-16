package mfx

import (
	"github.com/pkg/errors"
	"go.uber.org/config"
)

func ProvideConfig[TConfig any](configSectionName string) func(provider config.Provider) (TConfig, error) {
	return func(provider config.Provider) (TConfig, error) {
		var conf TConfig
		if err := provider.Get(configSectionName).Populate(&conf); err != nil {
			var empty TConfig
			return empty, errors.WithStack(err)
		}

		return conf, nil
	}
}
