package jobs

import (
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/handlers"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type NotificationEventReceiver struct {
	fx.In

	Receiver receiver.EventReceiver `name:"notification_event_receiver"`
}

type NotificationEventReaderJob struct {
	fx.Out

	Job job.BackgroundJob `name:"notification_event_reader_job"`
}

func provideNotificationEventReaderJob(
	rcv NotificationEventReceiver,
	handler handlers.NotificationEventHandler,
	logger *logrus.Logger,
) NotificationEventReaderJob {
	return NotificationEventReaderJob{
		Job: receiver.NewEventReceiverJob(rcv.Receiver, handler, logger),
	}
}
