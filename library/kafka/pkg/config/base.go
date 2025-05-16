package config

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

const (
	securityProtocol = "SASL_SSL"
	saslMechanisms   = "PLAIN"
)

type BaseKafkaUserConfig struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	SSLRootCAPath    string `yaml:"ssl_root_ca_path"`
}

func (config BaseKafkaUserConfig) ToKafkaConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"security.protocol": securityProtocol,
		"sasl.mechanisms":   saslMechanisms,
		"sasl.username":     config.Username,
		"sasl.password":     config.Password,
		"ssl.ca.location":   config.SSLRootCAPath,
	}
}
