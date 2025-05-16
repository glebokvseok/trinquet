package psql

import "time"

const configSectionName = "psql"

type databaseConfig struct {
	ConnectionString      string        `yaml:"connection_string"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
}
