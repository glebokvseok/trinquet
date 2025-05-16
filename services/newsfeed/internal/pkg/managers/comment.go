package managers

import (
	"context"
	"github.com/google/uuid"
	playersvc "github.com/move-mates/trinquet/library/playersvc/pkg"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
)

const (
	commentPerRequestCount = 20
)

type CommentManager interface {
	GetCommentSection(ctx context.Context, userID uuid.UUID, postID uuid.UUID, cursor int64) (models.CommentSection, error)
}

func ProvideCommentManager(
	postRepository repos.PostRepository,
	playerService playersvc.PlayerService,
) CommentManager {
	return &commentManager{
		postRepository: postRepository,
		playerService:  playerService,
	}
}

type commentManager struct {
	postRepository repos.PostRepository
	playerService  playersvc.PlayerService
}

func (mgr *commentManager) GetCommentSection(
	ctx context.Context,
	userID uuid.UUID,
	postID uuid.UUID,
	cursor int64,
) (models.CommentSection, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	comments, err := mgr.postRepository.GetPostComments(ctx, postID, cursor, commentPerRequestCount)
	if err != nil {
		return models.CommentSection{}, err
	}

	var userIDs []uuid.UUID
	mp := make(map[uuid.UUID]struct{})
	for i, comment := range comments {
		if _, ok := mp[comment.Author.ID]; !ok {
			userIDs = append(userIDs, comment.Author.ID)
			mp[comment.Author.ID] = struct{}{}
		}

		comments[i].SelfAuthored = comment.Author.ID == userID

		for j, reply := range comment.Replies {
			if _, ok := mp[reply.Author.ID]; !ok {
				userIDs = append(userIDs, reply.Author.ID)
				mp[reply.Author.ID] = struct{}{}
			}

			comment.Replies[j].SelfAuthored = reply.Author.ID == userID
		}
	}

	authorPreviews, err := mgr.playerService.GetPlayerPreviews(ctx, userIDs)
	if err != nil {
		return models.CommentSection{}, err
	}

	newCursor := cursor
	for i, comment := range comments {
		comments[i].Author.Complete(authorPreviews[comment.Author.ID])
		for j, reply := range comment.Replies {
			comment.Replies[j].Author.Complete(authorPreviews[reply.Author.ID])
		}

		newCursor = min(newCursor, comment.Timestamp)
	}

	return models.CommentSection{
		Comments:        comments,
		Cursor:          newCursor,
		HasMoreComments: len(comments) == commentPerRequestCount,
	}, nil
}
