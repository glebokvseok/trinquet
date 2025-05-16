package jobs

import (
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/handlers"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type NewsfeedEventReceiver struct {
	fx.In

	Receiver receiver.EventReceiver `name:"newsfeed_event_receiver"`
}

type NewsfeedEventReaderJob struct {
	fx.Out

	Job job.BackgroundJob `name:"newsfeed_event_reader_job"`
}

func provideNewsfeedEventReaderJob(
	rcv NewsfeedEventReceiver,
	handler handlers.NewsfeedEventHandler,
	logger *logrus.Logger,
) NewsfeedEventReaderJob {
	return NewsfeedEventReaderJob{
		Job: receiver.NewEventReceiverJob(rcv.Receiver, handler, logger),
	}
}
