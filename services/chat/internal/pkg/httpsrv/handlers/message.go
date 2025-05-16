package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/move-mates/trinquet/library/kafka/pkg/event/sender"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/events"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
	"time"
)

type MessageEventHandler struct {
	messageEventSender sender.EventSender
	validator          *validator.Validate
}

type MessageEventSender struct {
	fx.In

	Sender sender.EventSender `name:"message_event_sender"`
}

func provideMessageEventHandler(
	messageEventSender MessageEventSender,
	validator *validator.Validate,
) *MessageEventHandler {
	return &MessageEventHandler{
		messageEventSender: messageEventSender.Sender,
		validator:          validator,
	}
}

func (handler *MessageEventHandler) NewMessage(ctx echo.Context) error {
	return processEvent[events.NewMessageEvent](ctx, handler.messageEventSender, handler.validator)
}

func processEvent[TEvent any](
	ctx echo.Context,
	sender sender.EventSender,
	validator *validator.Validate,
) error {
	var ev TEvent
	var eventType events.MessageEventType

	switch any(ev).(type) {
	case events.NewMessageEvent:
		eventType = events.NewMessage
	default:
		return errors.New("unknown message event type")
	}

	err := ctx.Bind(&ev)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = validator.Struct(ev)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	rawEvent, err := json.Marshal(ev)
	if err != nil {
		return errors.WithStack(err)
	}

	wrapper := event.Wrapper{
		EventType: eventType,
		RawEvent:  rawEvent,
		UserID:    auth.GetUserID(ctx),
		RequestID: request.GetIDFromContext(ctx.Request().Context()),
		Timestamp: time.Now(),
	}

	err = sender.SendEvent(wrapper)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
