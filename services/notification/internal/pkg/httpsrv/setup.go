package httpsrv

import (
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	handlerspkg "github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv/handlers"
	middlewarepkg "github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv/middleware"
)

const (
	baseGroupPrefix = "/v0/notification-service"

	notificationGroupPrefix = "/notification"
)

func ProvideHttpServerSetupFunc(
	handlers handlerspkg.Params,
	middleware middlewarepkg.Params,
) httpsrv.HttpServerSetupFunc {
	return func(e *echo.Echo) {
		registerRoutes(e, handlers, middleware)
	}
}

func registerRoutes(
	e *echo.Echo,
	handlers handlerspkg.Params,
	middleware middlewarepkg.Params,
) {
	base := e.Group(baseGroupPrefix)
	base.Use(middleware.PanicHandler)
	base.Use(middleware.ErrorHandler)
	base.Use(middleware.CORS)
	base.GET(httpsrv.HealthCheckPath, httpsrv.HealthCheckFunc)
	base.Use(middleware.SignatureHandler)
	base.Use(middleware.RequestIDHandler)
	base.Use(middleware.UserAuthHandler)

	notification := base.Group(notificationGroupPrefix)
	notification.GET("/all", handlers.NotificationHandler.GetNotifications)
}
