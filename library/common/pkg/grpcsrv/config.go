package grpcsrv

const configSectionName string = "grpc_server"

type grpcServerConfig struct {
	Port uint16 `yaml:"port"`
}
