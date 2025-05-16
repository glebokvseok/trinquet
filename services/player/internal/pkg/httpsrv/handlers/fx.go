package handlers

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"handlers",
		fx.Provide(
			provideAvatarHandler,
			providePlayerHandler,
			providePlayerRelationHandler,
			provideRacquetMatchHandler,
			provideRacquetProfileHandler,
		),
	)
}

type Params struct {
	fx.In

	AvatarHandler         *AvatarHandler
	PlayerHandler         *PlayerHandler
	PlayerRelationHandler *PlayerRelationHandler
	RacquetMatchHandler   *RacquetMatchHandler
	RacquetProfileHandler *RacquetProfileHandler
}
