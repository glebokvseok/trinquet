package panicmw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/move-mates/trinquet/library/common/pkg/config"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type PanicHandlerMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"panic_handler"`
}

func providePanicHandlerMiddleware(appConfig appConfig, logger *logrus.Logger) PanicHandlerMiddleware {
	return PanicHandlerMiddleware{
		Middleware: middleware.RecoverWithConfig(
			middleware.RecoverConfig{
				LogErrorFunc: func(ctx echo.Context, err error, stack []byte) error {
					logger.
						WithContext(ctx.Request().Context()).
						Errorf("unhandeled panic occured: %s\n%+v", err.Error(), string(stack))

					message := "Internal Server Error"
					if appConfig.Mode != config.Prod {
						message = err.Error()
					}

					return ctx.JSON(
						http.StatusInternalServerError,
						errors.APIErrorResponse{
							Type:    errors.InternalServerError,
							Message: message,
						},
					)
				},
			},
		),
	}
}
