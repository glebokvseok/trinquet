package grpcsrv

import (
	"google.golang.org/grpc"
)

type GrpcServerSetupFunc func(server *grpc.Server)
