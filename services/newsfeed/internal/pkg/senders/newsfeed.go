package senders

import (
	"github.com/move-mates/trinquet/library/kafka/pkg/config"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/sender"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	NewsfeedEventSenderConfigSectionName = "newsfeed_event"

	newsfeedEventSenderName = "newsfeed-event"
)

type NewsfeedEventSenderConfig struct {
	config.EventSenderConfig `yaml:",inline"`
}

type NewsfeedEventSender struct {
	fx.Out

	Sender sender.EventSender `name:"newsfeed_event_sender"`
}

func ProvideNewsfeedEventSender(
	config NewsfeedEventSenderConfig,
	logger *logrus.Logger,
) (NewsfeedEventSender, error) {
	snd, err := sender.NewEventSender(config.EventSenderConfig, logger, newsfeedEventSenderName)

	return NewsfeedEventSender{
		Sender: snd,
	}, err
}
