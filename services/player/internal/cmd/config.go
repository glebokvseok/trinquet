package main

import "github.com/move-mates/trinquet/library/common/pkg/config"

func provideAppConfig() (config.AppConfig, error) {
	return config.ProvideAppConfig()
}
