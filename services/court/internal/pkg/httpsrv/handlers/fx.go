package handlers

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"handlers",
		fx.Provide(
			provideCourtHandler,
		),
	)
}

type Params struct {
	fx.In

	CourtHandler *CourtHandler
}
