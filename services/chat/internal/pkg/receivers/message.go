package receivers

import (
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	MessageEventReceiverConfigSectionName = "message_event"

	messageEventReceiverName = "message-event"
)

type MessageEventReceiverConfig struct {
	config.EventReceiverConfig `yaml:",inline"`
}

type MessageEventReceiver struct {
	fx.Out

	Receiver receiver.EventReceiver `name:"message_event_receiver"`
}

func ProvideMessageEventReceiver(
	config MessageEventReceiverConfig,
	logger *logrus.Logger,
) (MessageEventReceiver, error) {
	rcv, err := receiver.NewEventReceiver(config.EventReceiverConfig, logger, messageEventReceiverName)

	return MessageEventReceiver{
		Receiver: rcv,
	}, err
}
