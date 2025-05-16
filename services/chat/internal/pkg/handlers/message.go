package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/events"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	messageEventProcessingTimeout = time.Second * 5
)

type MessageEventHandler interface {
	event.Handler
}

type messageEventHandler struct {
	chatRepository    repos.ChatRepository
	messageRepository repos.MessageRepository
	logger            *logrus.Logger
}

func ProvideMessageEventHandler(
	chatRepository repos.ChatRepository,
	messageRepository repos.MessageRepository,
	logger *logrus.Logger,
) MessageEventHandler {
	return &messageEventHandler{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
		logger:            logger,
	}
}

func (handler *messageEventHandler) Handle(wrapper event.Wrapper) {
	ctx, cancel := context.WithTimeout(
		request.CreateContext(context.Background(), wrapper.RequestID),
		messageEventProcessingTimeout,
	)

	defer cancel()

	handler.logger.WithContext(ctx).Infof("start handling message event, type: %s", wrapper.EventType)

	var err error
	switch wrapper.EventType {
	case events.NewMessage:
		err = handler.handleNewMessageEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	default:
		err = errors.Errorf("unknown event type: %s", wrapper.EventType)
	}

	if err != nil {
		handler.logger.WithContext(ctx).Errorf("failed to handle message event, error: %+v", err)
	} else {
		handler.logger.WithContext(ctx).Infof("successfully handled message event")
	}
}

func (handler *messageEventHandler) handleNewMessageEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.NewMessageEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	chat, err := handler.chatRepository.GetInternalChat(ctx, userID, ev.Chat)
	if err != nil {
		return err
	}

	return handler.messageRepository.SaveMessage(ctx, userID, chat.ID, ev)
}
