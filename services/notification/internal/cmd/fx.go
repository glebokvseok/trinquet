package main

import (
	"context"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/library/common/pkg/logging"
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/handlers"
	setup "github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv"
	apihandlers "github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv/middleware"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/jobs"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/receivers"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

func App() *fx.App {
	return fx.New(
		fx.Provide(
			provideAppConfig,
		),
		fx.Options(
			apihandlers.Module(),
			httpsrv.Module(),
			jobs.Module(),
			logging.Module(),
			middleware.Module(),
			psql.Module(),
		),
		fx.Provide(
			setup.ProvideHttpServerSetupFunc,
			receivers.ProvideNotificationEventReceiver,
			mfx.ProvideConfig[receivers.NotificationEventReceiverConfig](receivers.NotificationEventReceiverConfigSectionName),
			handlers.ProvideNotificationEventHandler,
			managers.ProvideNotificationManager,
			repos.ProvideNotificationRepository,
		),
		fx.Invoke(
			func(jobs jobs.Params, logger *logrus.Logger) {
				ctx, cancel := context.WithCancel(context.Background())

				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

				go func() {
					<-sigChan
					cancel()
				}()

				for _, job := range jobs.List() {
					go func() {
						err := job.Run(ctx)
						if err != nil {
							logger.Fatalf("failed to start background job: %+v", err)
						}
					}()
				}
			},
			func(srv httpsrv.HttpServer) error {
				return srv.Run()
			},
		),
	)
}
