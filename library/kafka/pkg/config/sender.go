package config

type EventSenderConfig struct {
	BaseKafkaUserConfig `yaml:",inline"`
	Topic               string `yaml:"topic"`
}
