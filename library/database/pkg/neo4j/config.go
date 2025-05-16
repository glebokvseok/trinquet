package neo4j

import "time"

const configSectionName = "neo4j"

type SessionConfig struct {
	Database       string        `yaml:"database"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

type databaseConfig struct {
	URI                   string        `yaml:"uri"`
	Username              string        `yaml:"username"`
	Password              string        `yaml:"password"`
	SSLRootCert           string        `yaml:"ssl_root_cert"`
	MaxConnectionPoolSize int           `yaml:"max_connection_pool_size"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
}
