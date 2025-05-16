package grpcsrv

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"grpcsrv",
		fx.Provide(
			provideGrpcServer,
			mfx.ProvideConfig[grpcServerConfig](configSectionName),
		),
	)
}
