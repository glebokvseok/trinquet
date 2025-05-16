package main

import (
	"github.com/move-mates/trinquet/library/common/pkg/grpcsrv"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/library/common/pkg/logging"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"github.com/move-mates/trinquet/library/database/pkg/neo4j"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	grpcsetup "github.com/move-mates/trinquet/services/player/internal/pkg/grpcsrv"
	httpsetup "github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/middleware"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	vldconf "github.com/move-mates/trinquet/services/player/internal/pkg/validation"
	"github.com/sirupsen/logrus"
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
			grpcsrv.Module(),
			logging.Module(),
			middleware.Module(),
			neo4j.Module(),
			psql.Module(),
			s3.Module(),
			validation.Module(),
		),
		fx.Provide(
			httpsetup.ProvideHttpServerSetupFunc,
			grpcsetup.ProvideGrpcServerSetupFunc,
			managers.ProvideAvatarManager,
			managers.ProvidePlayerManager,
			managers.ProvideRelationManager,
			managers.ProvideRacquetProfileManager,
			managers.ProvideRacquetMatchManager,
			managers.ProvideAchievementManager,
			repos.ProvideTransactionManager,
			repos.ProvideAvatarRepository,
			repos.ProvidePlayerRepository,
			repos.ProvideRelationRepository,
			repos.ProvideRacquetProfileRepository,
			repos.ProvideAchievementRepository,
			repos.ProvideRacquetMatchRepository,
			vldconf.ProvideValidatorConfigFunc,
		),
		fx.Invoke(
			func(srv grpcsrv.GrpcServer, logger *logrus.Logger) error {
				go func() {
					err := srv.Run()
					if err != nil {
						logger.Fatalf("failed to start grpc server: %+v", err)
					}
				}()

				return nil
			},
			func(srv httpsrv.HttpServer) error {
				return srv.Run()
			},
		),
	)
}
