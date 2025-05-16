package logging

const configSectionName = "logging"

type loggerConfig struct {
	LogLevel string `yaml:"log_level"`
}
