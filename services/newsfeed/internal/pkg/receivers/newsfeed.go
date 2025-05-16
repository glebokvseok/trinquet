package receivers

import (
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/receiver"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	NewsfeedEventReceiverConfigSectionName = "newsfeed_event"

	newsfeedEventReceiverName = "newsfeed-event"
)

type NewsfeedEventReceiverConfig struct {
	config.EventReceiverConfig `yaml:",inline"`
}

type NewsfeedEventReceiver struct {
	fx.Out

	Receiver receiver.EventReceiver `name:"newsfeed_event_receiver"`
}

func ProvideNewsfeedEventReceiver(
	config NewsfeedEventReceiverConfig,
	logger *logrus.Logger,
) (NewsfeedEventReceiver, error) {
	rcv, err := receiver.NewEventReceiver(config.EventReceiverConfig, logger, newsfeedEventReceiverName)

	return NewsfeedEventReceiver{
		Receiver: rcv,
	}, err
}
