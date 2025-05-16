package handlers

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"handler",
		fx.Provide(
			provideNotificationHandler,
		),
	)
}

type Params struct {
	fx.In

	NotificationHandler *NotificationHandler
}
