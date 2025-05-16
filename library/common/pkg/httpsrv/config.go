package httpsrv

const configSectionName string = "http_server"

type httpServerConfig struct {
	Port uint16 `yaml:"port"`
}
