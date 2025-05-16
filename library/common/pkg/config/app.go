package config

import (
	"github.com/pkg/errors"
	uconf "go.uber.org/config"
	"go.uber.org/fx"
)

type AppConfig struct {
	fx.Out

	Provider uconf.Provider
}

func ProvideAppConfig() (AppConfig, error) {
	err := PopulateEnvVars()
	if err != nil {
		return AppConfig{}, errors.WithStack(err)
	}

	envProvider, err := GetAppEnvConfigProvider()
	if err != nil {
		return AppConfig{}, errors.WithStack(err)
	}

	modeProvider, err := GetAppModeConfigProvider()
	if err != nil {
		return AppConfig{}, errors.WithStack(err)
	}

	provider, err := uconf.NewProviderGroup("config", modeProvider, envProvider)
	if err != nil {
		return AppConfig{}, errors.WithStack(err)
	}

	return AppConfig{
		Provider: provider,
	}, nil
}
