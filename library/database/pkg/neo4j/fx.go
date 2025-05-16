package neo4j

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"neo4j",
		fx.Provide(
			provideDriver,
			mfx.ProvideConfig[databaseConfig](configSectionName),
			mfx.ProvideConfig[SessionConfig](configSectionName),
		),
	)
}
