package grpcsrv

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type GrpcServer interface {
	Run() error
}

type grpcServer struct {
	server   *grpc.Server
	listener net.Listener
}

func provideGrpcServer(
	logger *logrus.Logger,
	config grpcServerConfig,
	setupServer GrpcServerSetupFunc,
) (GrpcServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(newPanicInterceptor(logger)),
	)

	setupServer(server)

	return &grpcServer{
		server:   server,
		listener: listener,
	}, nil
}

func (s *grpcServer) Run() error {
	err := s.server.Serve(s.listener)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func newPanicInterceptor(
	logger *logrus.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("unhandeled panic occured in method %s: %+v", info.FullMethod, r)

				err = status.Errorf(codes.Internal, "Internal Server Error")
			}
		}()

		return handler(ctx, req)
	}
}
