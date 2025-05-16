package playersvc

import "time"

const (
	configSectionName = "player_service"
)

type playerServiceConfig struct {
	Host           string        `yaml:"host"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}
