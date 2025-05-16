package main

import (
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/library/common/pkg/logging"
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/generators"
	setup "github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/middleware"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/managers"
	vldconf "github.com/move-mates/trinquet/services/auth/internal/pkg/validation"
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
			validation.Module(),
		),
		fx.Provide(
			setup.ProvideHttpServerSetupFunc,
			managers.ProvideAuthenticationManager,
			managers.ProvideRegistrationManager,
			generators.ProvideJWTGenerator,
			mfx.ProvideConfig[generators.JWTConfig](generators.JWTConfigSectionName),
			repos.ProvideTransactionManager,
			repos.ProvideUserRepository,
			vldconf.ProvideValidatorConfigFunc,
		),
		fx.Invoke(
			func(srv httpsrv.HttpServer) error {
				return srv.Run()
			},
		),
	)
}
