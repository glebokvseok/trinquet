package reqidmw

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"request_id_middleware",
		fx.Provide(
			provideRequestIDMiddleware,
		),
	)
}
