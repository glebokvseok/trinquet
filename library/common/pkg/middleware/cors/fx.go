package corsmw

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"cors_middleware",
		fx.Provide(
			provideCORSMiddleware,
		),
	)
}
