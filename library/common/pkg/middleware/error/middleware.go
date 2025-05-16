package errmw

import (
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/config"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	goerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type ErrorHandlerMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"error_handler"`
}

func provideErrorHandlerMiddleware(
	config appConfig,
	logger *logrus.Logger,
) ErrorHandlerMiddleware {
	return ErrorHandlerMiddleware{
		Middleware: errorHandler(config, logger),
	}
}

func errorHandler(appConfig appConfig, logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			if err != nil {
				if apiErr := (*errors.APIError)(nil); goerrors.As(err, &apiErr) {
					return ctx.JSON(
						apiErr.HttpCode,
						errors.APIErrorResponse{
							Type:    apiErr.Type,
							Message: apiErr.Message,
						},
					)
				} else if httpErr := (*echo.HTTPError)(nil); goerrors.As(err, &httpErr) {
					return httpErr
				} else {
					logger.
						WithContext(ctx.Request().Context()).
						Errorf("unhandeled error occured: %+v", err)

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
				}
			}

			return nil
		}
	}
}
