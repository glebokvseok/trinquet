package httpsrv

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HttpServer interface {
	Run() error
}

type httpServer struct {
	echo   *echo.Echo
	config httpServerConfig
}

func provideHttpServer(
	config httpServerConfig,
	setupServer HttpServerSetupFunc,
) HttpServer {
	e := echo.New()

	setupServer(e)

	return &httpServer{
		echo:   e,
		config: config,
	}
}

func (s *httpServer) Run() error {
	err := s.echo.Start(fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
