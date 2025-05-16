package config

import (
	"github.com/pkg/errors"
	"go.uber.org/config"
	"os"
)

const (
	Dev  = "development"
	Test = "testing"
	Prod = "production"
)

const (
	devConfigFilePath  = "config/config.dev.yaml"
	testConfigFilePath = "config/config.test.yaml"
	prodConfigFilePath = "config/config.prod.yaml"
)

func GetAppModeConfigProvider() (*config.YAML, error) {
	configFilePath, err := getAppModeConfigFilePath()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	yaml, err := config.NewYAML(config.File(configFilePath))
	return yaml, errors.WithStack(err)
}

func getAppModeConfigFilePath() (string, error) {
	switch mode := os.Getenv("APP_MODE"); mode {
	case Dev:
		return devConfigFilePath, nil
	case Test:
		return testConfigFilePath, nil
	case Prod:
		return prodConfigFilePath, nil
	default:
		return "", errors.Errorf("unknown application mode: %s", mode)
	}
}
