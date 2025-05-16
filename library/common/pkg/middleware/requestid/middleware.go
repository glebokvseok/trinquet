package reqidmw

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/constants"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type RequestIDMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"request_id_handler"`
}

func provideRequestIDMiddleware(
	logger *logrus.Logger,
) RequestIDMiddleware {
	return RequestIDMiddleware{
		Middleware: func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				var requestID uuid.UUID
				requestIDHeader := ctx.Request().Header.Get(constants.RequestIDHeader)
				if extensions.IsNotEmpty(requestIDHeader) {
					var err error
					requestID, err = uuid.Parse(requestIDHeader)

					if err != nil {
						return NewInvalidRequestIdFormatError()
					}
				} else {
					requestID = uuid.New()
				}

				newCtx := request.CreateContext(ctx.Request().Context(), requestID)
				ctx.SetRequest(ctx.Request().WithContext(newCtx))

				logger.WithContext(newCtx).Infof("start handling request with path: %s", ctx.Path())

				return next(ctx)
			}
		},
	}
}
