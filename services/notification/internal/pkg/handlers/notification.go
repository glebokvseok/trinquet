package handlers

import (
	"context"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	notificationEventProcessingTimeout = time.Second * 5
)

type NotificationEventHandler interface {
	event.Handler
}

type notificationEventHandler struct {
	logger *logrus.Logger
}

func ProvideNotificationEventHandler(
	logger *logrus.Logger,
) NotificationEventHandler {
	return &notificationEventHandler{
		logger: logger,
	}
}

func (handler *notificationEventHandler) Handle(wrapper event.Wrapper) {
	ctx, cancel := context.WithTimeout(
		request.CreateContext(context.Background(), wrapper.RequestID),
		notificationEventProcessingTimeout,
	)

	defer cancel()

	handler.logger.WithContext(ctx).Infof("start handling notification event, type: %s", wrapper.EventType)
}
