package handlers

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"handlers",
		fx.Provide(
			provideAuthenticationHandler,
			provideRegistrationHandler,
		),
	)
}

type Params struct {
	fx.In

	AuthenticationHandler *AuthenticationHandler
	RegistrationHandler   *RegistrationHandler
}
