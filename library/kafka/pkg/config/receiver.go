package config

type EventReceiverConfig struct {
	BaseKafkaUserConfig `yaml:",inline"`
	Topic               string `yaml:"topic"`
	ConsumerGroup       string `yaml:"consumer_group"`
}
