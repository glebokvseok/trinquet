package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/move-mates/trinquet/library/common/pkg/result"
	"github.com/move-mates/trinquet/library/kafka/pkg/event"
	"github.com/move-mates/trinquet/library/playersvc/pkg"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/events"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/managers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	newsfeedEventProcessingTimeout = time.Second * 5
)

type NewsfeedEventHandler interface {
	event.Handler
}

type newsfeedEventHandler struct {
	playerService   playersvc.PlayerService
	postManager     managers.PostManager
	newsfeedManager managers.NewsfeedManager
	logger          *logrus.Logger
}

func ProvideNewsfeedEventHandler(
	playerService playersvc.PlayerService,
	postManager managers.PostManager,
	newsfeedManager managers.NewsfeedManager,
	logger *logrus.Logger,
) NewsfeedEventHandler {
	return &newsfeedEventHandler{
		playerService:   playerService,
		postManager:     postManager,
		newsfeedManager: newsfeedManager,
		logger:          logger,
	}
}

func (handler *newsfeedEventHandler) Handle(wrapper event.Wrapper) {
	ctx, cancel := context.WithTimeout(
		request.CreateContext(context.Background(), wrapper.RequestID),
		newsfeedEventProcessingTimeout,
	)

	defer cancel()

	handler.logger.WithContext(ctx).Infof("start handling newsfeed event, type: %s", wrapper.EventType)

	var err error
	switch wrapper.EventType {
	case events.CreatePost:
		err = handler.handleCreatePostEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	case events.LikePost:
		err = handler.handleLikePostEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	case events.UnlikePost:
		err = handler.handleUnlikePostEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	case events.CommentPost:
		err = handler.handleCommentPostEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	case events.ReplyPostComment:
		err = handler.handleReplyPostCommentEvent(ctx, wrapper.UserID, wrapper.RawEvent)
	default:
		err = errors.Errorf("unknown event type: %s", wrapper.EventType)
	}

	if err != nil {
		handler.logger.WithContext(ctx).Errorf("failed to handle newsfeed event, error: %+v", err)
	} else {
		handler.logger.WithContext(ctx).Infof("successfully handled newsfeed event")
	}
}

func (handler *newsfeedEventHandler) handleCreatePostEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.CreatePostEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	getFollowersCh := make(chan result.Result[[]uuid.UUID])
	go handler.getFollowers(ctx, userID, getFollowersCh)

	postID, err := handler.postManager.CreatePost(ctx, userID, ev)
	if err != nil {
		return err
	}

	followersRes := <-getFollowersCh
	if followersRes.Err != nil {
		return followersRes.Err
	}

	return handler.newsfeedManager.AddPostToNewsfeeds(ctx, followersRes.Data, postID, ev.CreatedOn)
}

func (handler *newsfeedEventHandler) handleLikePostEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.LikePostEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	return handler.postManager.LikePost(ctx, userID, ev)
}

func (handler *newsfeedEventHandler) handleUnlikePostEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.UnlikePostEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	return handler.postManager.UnlikePost(ctx, userID, ev)
}

func (handler *newsfeedEventHandler) handleCommentPostEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.CommentPostEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	return handler.postManager.CommentPost(ctx, userID, ev)
}

func (handler *newsfeedEventHandler) handleReplyPostCommentEvent(ctx context.Context, userID uuid.UUID, rawEvent []byte) error {
	var ev events.ReplyPostCommentEvent
	if err := json.Unmarshal(rawEvent, &ev); err != nil {
		return errors.WithStack(err)
	}

	return handler.postManager.ReplyPostComment(ctx, userID, ev)
}

func (handler *newsfeedEventHandler) getFollowers(ctx context.Context, userID uuid.UUID, ch chan<- result.Result[[]uuid.UUID]) {
	followers, err := handler.playerService.GetAllFollowers(ctx, userID)
	if err != nil {
		ch <- result.Result[[]uuid.UUID]{Err: err}

		return
	}

	ch <- result.Result[[]uuid.UUID]{Data: followers}
}
