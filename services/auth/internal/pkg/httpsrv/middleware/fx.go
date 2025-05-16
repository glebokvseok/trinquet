package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/middleware/cors"
	"github.com/move-mates/trinquet/library/common/pkg/middleware/error"
	"github.com/move-mates/trinquet/library/common/pkg/middleware/panic"
	reqidmw "github.com/move-mates/trinquet/library/common/pkg/middleware/requestid"
	"github.com/move-mates/trinquet/library/common/pkg/middleware/signature"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"middleware",
		fx.Options(
			corsmw.Module(),
			errmw.Module(),
			panicmw.Module(),
			reqidmw.Module(),
			signmw.Module(),
		),
		fx.Provide(
			provideUserRefreshMiddleware,
		),
	)
}

type Params struct {
	fx.In

	CORS               echo.MiddlewareFunc `name:"cors"`
	ErrorHandler       echo.MiddlewareFunc `name:"error_handler"`
	PanicHandler       echo.MiddlewareFunc `name:"panic_handler"`
	RequestIDHandler   echo.MiddlewareFunc `name:"request_id_handler"`
	SignatureHandler   echo.MiddlewareFunc `name:"signature_handler"`
	UserRefreshHandler echo.MiddlewareFunc `name:"user_refresh_handler"`
}
