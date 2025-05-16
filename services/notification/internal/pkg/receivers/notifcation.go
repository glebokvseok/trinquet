package receivers

import (
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	NotificationEventReceiverConfigSectionName = "notification_event"

	notificationEventReceiverName = "notification-event"
)

type NotificationEventReceiverConfig struct {
	config.EventReceiverConfig `yaml:",inline"`
}

type NotificationEventReceiver struct {
	fx.Out

	Receiver receiver.EventReceiver `name:"notification_event_receiver"`
}

func ProvideNotificationEventReceiver(
	config NotificationEventReceiverConfig,
	logger *logrus.Logger,
) (NotificationEventReceiver, error) {
	rcv, err := receiver.NewEventReceiver(config.EventReceiverConfig, logger, notificationEventReceiverName)

	return NotificationEventReceiver{
		Receiver: rcv,
	}, err
}
