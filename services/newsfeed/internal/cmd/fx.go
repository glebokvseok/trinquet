package main

import (
	"context"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/library/common/pkg/logging"
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"github.com/move-mates/trinquet/library/common/pkg/validation"
	"github.com/move-mates/trinquet/library/database/pkg/mongo"
	playersvc "github.com/move-mates/trinquet/library/playersvc/pkg"
	"github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/handlers"
	setup "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/httpsrv"
	apihandlers "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/httpsrv/middleware"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/jobs"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/receivers"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/senders"
	vldconf "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/validation"
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
			mongo.Module(),
			playersvc.Module(),
			s3.Module(),
			validation.Module(),
		),
		fx.Provide(
			setup.ProvideHttpServerSetupFunc,
			senders.ProvideNewsfeedEventSender,
			mfx.ProvideConfig[senders.NewsfeedEventSenderConfig](senders.NewsfeedEventSenderConfigSectionName),
			receivers.ProvideNewsfeedEventReceiver,
			mfx.ProvideConfig[receivers.NewsfeedEventReceiverConfig](receivers.NewsfeedEventReceiverConfigSectionName),
			handlers.ProvideNewsfeedEventHandler,
			managers.ProvideCommentManager,
			managers.ProvidePostManager,
			managers.ProvideNewsfeedManager,
			managers.ProvideMediaManager,
			repos.ProvidePostRepository,
			repos.ProvideNewsfeedRepository,
			repos.ProvideMediaRepository,
			vldconf.ProvideValidatorConfigFunc,
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
