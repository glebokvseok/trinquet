package mongo

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"mongo",
		fx.Provide(
			provideClient,
			mfx.ProvideConfig[databaseConfig](configSectionName),
			mfx.ProvideConfig[RequestConfig](configSectionName),
		),
	)
}
