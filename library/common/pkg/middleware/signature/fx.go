package signmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"signature_handler_middleware",
		fx.Provide(
			provideSignatureHandlerMiddleware,
			mfx.ProvideConfig[signatureConfig](configSectionName),
		),
	)
}
