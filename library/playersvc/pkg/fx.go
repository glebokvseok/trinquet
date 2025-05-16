package playersvc

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"playersvc",
		fx.Provide(
			providePlayerService,
			mfx.ProvideConfig[playerServiceConfig](configSectionName),
		),
	)
}
