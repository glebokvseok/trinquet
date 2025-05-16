package jobs

import (
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/handlers"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type MessageEventReceiver struct {
	fx.In

	Receiver receiver.EventReceiver `name:"message_event_receiver"`
}

type MessageEventReaderJob struct {
	fx.Out

	Job job.BackgroundJob `name:"message_event_reader_job"`
}

func provideMessageEventReaderJob(
	rcv MessageEventReceiver,
	handler handlers.MessageEventHandler,
	logger *logrus.Logger,
) MessageEventReaderJob {
	return MessageEventReaderJob{
		Job: receiver.NewEventReceiverJob(rcv.Receiver, handler, logger),
	}
}
