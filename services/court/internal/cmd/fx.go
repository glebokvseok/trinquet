package main

import (
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/library/common/pkg/logging"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/court/internal/pkg/database/repos"
	httpsetup "github.com/move-mates/trinquet/services/court/internal/pkg/httpsrv"
	"github.com/move-mates/trinquet/services/court/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/court/internal/pkg/httpsrv/middleware"
	"github.com/move-mates/trinquet/services/court/internal/pkg/managers"
	vldconf "github.com/move-mates/trinquet/services/court/internal/pkg/validation"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		fx.Provide(
			provideAppConfig,
		),
		fx.Options(
			handlers.Module(),
			httpsrv.Module(),
			logging.Module(),
			middleware.Module(),
			psql.Module(),
			s3.Module(),
			validation.Module(),
		),
		fx.Provide(
			httpsetup.ProvideHttpServerSetupFunc,
			managers.ProvideCourtManager,
			repos.ProvideCourtRepository,
			vldconf.ProvideValidatorConfigFunc,
		),
		fx.Invoke(
			func(srv httpsrv.HttpServer) error {
				return srv.Run()
			},
		),
	)
}
