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
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/events"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
	"time"
)

type PostEventHandler struct {
	newsfeedEventSender sender.EventSender
	validator           *validator.Validate
}

type NewsfeedEventSender struct {
	fx.In

	Sender sender.EventSender `name:"newsfeed_event_sender"`
}

func providePostEventHandler(
	newsfeedEventSender NewsfeedEventSender,
	validator *validator.Validate,
) *PostEventHandler {
	return &PostEventHandler{
		newsfeedEventSender: newsfeedEventSender.Sender,
		validator:           validator,
	}
}

func (handler *PostEventHandler) CreatePost(ctx echo.Context) error {
	return processEvent[events.CreatePostEvent](ctx, handler.newsfeedEventSender, handler.validator)
}

func (handler *PostEventHandler) LikePost(ctx echo.Context) error {
	return processEvent[events.LikePostEvent](ctx, handler.newsfeedEventSender, handler.validator)
}

func (handler *PostEventHandler) UnlikePost(ctx echo.Context) error {
	return processEvent[events.UnlikePostEvent](ctx, handler.newsfeedEventSender, handler.validator)
}

func (handler *PostEventHandler) CommentPost(ctx echo.Context) error {
	return processEvent[events.CommentPostEvent](ctx, handler.newsfeedEventSender, handler.validator)
}

func (handler *PostEventHandler) ReplyPostComment(ctx echo.Context) error {
	return processEvent[events.ReplyPostCommentEvent](ctx, handler.newsfeedEventSender, handler.validator)
}

func processEvent[TEvent any](
	ctx echo.Context,
	sender sender.EventSender,
	validator *validator.Validate,
) error {
	var ev TEvent
	var eventType events.PostEventType

	switch any(ev).(type) {
	case events.CreatePostEvent:
		eventType = events.CreatePost
	case events.LikePostEvent:
		eventType = events.LikePost
	case events.UnlikePostEvent:
		eventType = events.UnlikePost
	case events.CommentPostEvent:
		eventType = events.CommentPost
	case events.ReplyPostCommentEvent:
		eventType = events.ReplyPostComment
	default:
		return errors.New("unknown post event type")
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
