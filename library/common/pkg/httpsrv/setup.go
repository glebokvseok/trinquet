package httpsrv

import "github.com/labstack/echo/v4"

type HttpServerSetupFunc func(echo *echo.Echo)
