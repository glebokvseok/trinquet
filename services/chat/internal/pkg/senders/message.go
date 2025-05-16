package senders

import (
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/sender"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	MessageEventSenderConfigSectionName = "message_event"

	messageEventSenderName = "message-event"
)

type MessageEventSenderConfig struct {
	config.EventSenderConfig `yaml:",inline"`
}

type MessageEventSender struct {
	fx.Out

	Sender sender.EventSender `name:"message_event_sender"`
}

func ProvideMessageEventSender(
	config MessageEventSenderConfig,
	logger *logrus.Logger,
) (MessageEventSender, error) {
	snd, err := sender.NewEventSender(config.EventSenderConfig, logger, messageEventSenderName)

	return MessageEventSender{
		Sender: snd,
	}, err
}
