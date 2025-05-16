package receiver

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/move-mates/trinquet/library/common/pkg/job"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type receiverJob struct {
	receiver EventReceiver
	handler  event.Handler
	logger   *logrus.Logger
}

func NewEventReceiverJob(
	receiver EventReceiver,
	handler event.Handler,
	logger *logrus.Logger,
) job.BackgroundJob {
	return &receiverJob{
		receiver: receiver,
		handler:  handler,
		logger:   logger,
	}
}

func (job *receiverJob) Run(ctx context.Context) error {
	job.logger.Infof("%s receiver job started", job.receiver.Name())

	for {
		select {
		case <-ctx.Done():
			job.logger.Warnf("%s receiver job stopped", job.receiver.Name())

			return nil
		default:
			job.receiveEvent()
		}
	}
}

func (job *receiverJob) receiveEvent() {
	defer func() {
		if r := recover(); r != nil {
			job.logger.Errorf("panic occured in %s receiver job: %+v", job.receiver.Name(), r)
		}
	}()

	ev, err := job.receiver.ReceiveEvent(-1) // -1 for indefinite wait

	if err != nil {
		if kafkaErr := (kafka.Error{}); errors.As(err, &kafkaErr) {
			if kafkaErr.IsTimeout() {
				return
			}
		}

		job.logger.Errorf("error occured in %s receiver job: %+v", job.receiver.Name(), err)

		return
	}

	job.handler.Handle(ev)
}
