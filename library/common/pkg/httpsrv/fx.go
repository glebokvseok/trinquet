package httpsrv

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"httpsrv",
		fx.Provide(
			provideHttpServer,
			mfx.ProvideConfig[httpServerConfig](configSectionName),
		),
	)
}
