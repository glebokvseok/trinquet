package panicmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"panic_handler_middleware",
		fx.Provide(
			providePanicHandlerMiddleware,
			mfx.ProvideConfig[appConfig](configSectionName),
		),
	)
}
