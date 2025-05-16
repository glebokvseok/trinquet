package jobs

import (
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"jobs",
		fx.Provide(
			provideNotificationEventReaderJob,
		),
	)
}

type Params struct {
	fx.In

	NotificationEventReaderJob job.BackgroundJob `name:"notification_event_reader_job"`
}

func (params Params) List() []job.BackgroundJob {
	return []job.BackgroundJob{
		params.NotificationEventReaderJob,
	}
}
