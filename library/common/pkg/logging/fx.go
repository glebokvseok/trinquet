package logging

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"logging",
		fx.Provide(
			provideLogger,
			mfx.ProvideConfig[loggerConfig](configSectionName),
		),
	)
}
