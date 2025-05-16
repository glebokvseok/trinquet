package validation

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"validation",
		fx.Provide(
			provideValidator,
		),
	)
}
