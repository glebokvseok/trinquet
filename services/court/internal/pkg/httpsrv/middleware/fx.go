package middleware

import (
	"github.com/labstack/echo/v4"
	authmw "github.com/move-mates/trinquet/library/auth/pkg/middleware"
	corsmw "github.com/move-mates/trinquet/library/common/pkg/middleware/cors"
	errmw "github.com/move-mates/trinquet/library/common/pkg/middleware/error"
	panicmw "github.com/move-mates/trinquet/library/common/pkg/middleware/panic"
	reqidmw "github.com/move-mates/trinquet/library/common/pkg/middleware/requestid"
	signmw "github.com/move-mates/trinquet/library/common/pkg/middleware/signature"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"middleware",
		fx.Options(
			authmw.Module(),
			corsmw.Module(),
			errmw.Module(),
			panicmw.Module(),
			reqidmw.Module(),
			signmw.Module(),
		),
	)
}

type Params struct {
	fx.In

	CORS             echo.MiddlewareFunc `name:"cors"`
	ErrorHandler     echo.MiddlewareFunc `name:"error_handler"`
	PanicHandler     echo.MiddlewareFunc `name:"panic_handler"`
	RequestIDHandler echo.MiddlewareFunc `name:"request_id_handler"`
	SignatureHandler echo.MiddlewareFunc `name:"signature_handler"`
	UserAuthHandler  echo.MiddlewareFunc `name:"user_auth_handler"`
}
