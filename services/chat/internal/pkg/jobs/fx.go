package jobs

import (
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"jobs",
		fx.Provide(
			provideMessageEventReaderJob,
		),
	)
}

type Params struct {
	fx.In

	MessageEventReaderJob job.BackgroundJob `name:"message_event_reader_job"`
}

func (params Params) List() []job.BackgroundJob {
	return []job.BackgroundJob{
		params.MessageEventReaderJob,
	}
}
