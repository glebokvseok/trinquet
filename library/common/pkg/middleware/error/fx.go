package errmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"error_handler_middleware",
		fx.Provide(
			provideErrorHandlerMiddleware,
			mfx.ProvideConfig[appConfig](configSectionName),
		),
	)
}
