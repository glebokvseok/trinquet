package logging

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func provideLogger(config loggerConfig) (*logrus.Logger, error) {
	logger := logrus.New()

	err := configureLogger(logger, config)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func configureLogger(logger *logrus.Logger, config loggerConfig) error {
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return errors.WithStack(err)
	}

	logger.SetLevel(logLevel)

	logger.AddHook(&ContextHook{})

	return nil
}
