package mongo

import "time"

const configSectionName = "mongo"

type RequestConfig struct {
	Database       string        `yaml:"database"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

type databaseConfig struct {
	ConnectionString      string        `yaml:"connection_string"`
	MaxConnectionPoolSize uint64        `yaml:"max_connection_pool_size"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
}
