package authmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"user_auth_middleware",
		fx.Provide(
			provideUserAuthMiddleware,
			mfx.ProvideConfig[JWTConfig](configSectionName),
		),
	)
}
