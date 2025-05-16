package psql

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"psql",
		fx.Provide(
			provideDatabase,
			mfx.ProvideConfig[databaseConfig](configSectionName),
		),
	)
}
