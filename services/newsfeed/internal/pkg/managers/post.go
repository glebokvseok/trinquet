package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/events"
)

type PostManager interface {
	CreatePost(ctx context.Context, userID uuid.UUID, event events.CreatePostEvent) (postID uuid.UUID, err error)
	LikePost(ctx context.Context, userID uuid.UUID, event events.LikePostEvent) error
	UnlikePost(ctx context.Context, userID uuid.UUID, event events.UnlikePostEvent) error
	CommentPost(ctx context.Context, userID uuid.UUID, event events.CommentPostEvent) error
	ReplyPostComment(ctx context.Context, userID uuid.UUID, event events.ReplyPostCommentEvent) error
}

func ProvidePostManager(
	postRepository repos.PostRepository,
) PostManager {
	return &postManager{
		postRepository: postRepository,
	}
}

type postManager struct {
	postRepository repos.PostRepository
}

func (mgr *postManager) CreatePost(ctx context.Context, userID uuid.UUID, event events.CreatePostEvent) (uuid.UUID, error) {
	postID := uuid.New()

	err := mgr.postRepository.CreatePost(ctx, userID, postID, event)
	if err != nil {
		return uuid.Nil, err
	}

	return postID, nil
}

func (mgr *postManager) LikePost(ctx context.Context, userID uuid.UUID, event events.LikePostEvent) error {
	return mgr.postRepository.LikePost(ctx, userID, event)
}

func (mgr *postManager) UnlikePost(ctx context.Context, userID uuid.UUID, event events.UnlikePostEvent) error {
	return mgr.postRepository.UnlikePost(ctx, userID, event)
}

func (mgr *postManager) CommentPost(ctx context.Context, userID uuid.UUID, event events.CommentPostEvent) error {
	return mgr.postRepository.CommentPost(ctx, userID, event)
}

func (mgr *postManager) ReplyPostComment(ctx context.Context, userID uuid.UUID, event events.ReplyPostCommentEvent) error {
	return mgr.postRepository.ReplyPostComment(ctx, userID, event)
}
