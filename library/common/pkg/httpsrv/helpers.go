package httpsrv

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	HealthCheckPath = "/health-check"
)

func HealthCheckFunc(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
