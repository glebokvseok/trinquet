package httpsrv

import (
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/middleware"
)

const (
	baseGroupPrefix = "/v0/auth-service"

	userGroupPrefix = "/user"
)

func ProvideHttpServerSetupFunc(
	handlers handlers.Params,
	middleware middleware.Params,
) httpsrv.HttpServerSetupFunc {
	return func(e *echo.Echo) {
		registerRoutes(e, handlers, middleware)
	}
}

func registerRoutes(
	e *echo.Echo,
	handlers handlers.Params,
	middleware middleware.Params,
) {
	base := e.Group(baseGroupPrefix)
	base.Use(middleware.PanicHandler)
	base.Use(middleware.ErrorHandler)
	base.Use(middleware.CORS)
	base.GET(httpsrv.HealthCheckPath, httpsrv.HealthCheckFunc)
	base.Use(middleware.SignatureHandler)
	base.Use(middleware.RequestIDHandler)

	user := base.Group(userGroupPrefix)
	user.POST("/auth", handlers.AuthenticationHandler.AuthenticateUser)
	user.POST("/register", handlers.RegistrationHandler.RegisterUser)
	user.Use(middleware.UserRefreshHandler)
	user.GET("/refresh", handlers.AuthenticationHandler.RefreshUser)
}
