package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/result"
	playersvc "github.com/move-mates/trinquet/library/playersvc/pkg"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/helpers/media"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	postPerRequestCount = 20
)

type NewsfeedManager interface {
	AddPostToNewsfeeds(ctx context.Context, userIDs []uuid.UUID, postID uuid.UUID, postCreatedOn time.Time) error
	GetNewsfeed(ctx context.Context, userID uuid.UUID, cursor int64) (models.Newsfeed, error)
}

func ProvideNewsfeedManager(
	newsfeedRepository repos.NewsfeedRepository,
	postRepository repos.PostRepository,
	playerService playersvc.PlayerService,
	linkGenerator s3.LinkGenerator,
	logger *logrus.Logger,
) NewsfeedManager {
	return &newsfeedManager{
		newsfeedRepository: newsfeedRepository,
		postRepository:     postRepository,
		playerService:      playerService,
		linkGenerator:      linkGenerator,
		logger:             logger,
	}
}

type newsfeedManager struct {
	newsfeedRepository repos.NewsfeedRepository
	postRepository     repos.PostRepository
	playerService      playersvc.PlayerService
	linkGenerator      s3.LinkGenerator
	logger             *logrus.Logger
}

func (mgr *newsfeedManager) AddPostToNewsfeeds(
	ctx context.Context,
	userIDs []uuid.UUID,
	postID uuid.UUID,
	postCreatedOn time.Time,
) error {
	addedCount, err := mgr.newsfeedRepository.AddPostToNewsfeeds(ctx, userIDs, postID, postCreatedOn)
	if err != nil {
		return err
	}

	if addedCount != int64(len(userIDs)) {
		mgr.logger.WithContext(ctx).Warnf("post with was not added to all users newsfeeds, id: %s", postID)
	}

	return nil
}

func (mgr *newsfeedManager) GetNewsfeed(
	ctx context.Context,
	userID uuid.UUID,
	cursor int64,
) (models.Newsfeed, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	posts, err := mgr.newsfeedRepository.GetNewsfeedPosts(ctx, userID, cursor, postPerRequestCount)
	if err != nil {
		return models.Newsfeed{}, err
	}

	var userIDs []uuid.UUID
	mp := make(map[uuid.UUID]struct{})
	postIds := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIds[i] = post.ID
		if _, ok := mp[post.Author.ID]; !ok {
			userIDs = append(userIDs, post.Author.ID)
			mp[post.Author.ID] = struct{}{}
		}
	}

	getLikeCountCh := make(chan result.Result[[]int64])
	go mgr.getPostsLikeCount(ctx, postIds, getLikeCountCh)

	getCommentCountCh := make(chan result.Result[[]int64])
	go mgr.getPostsCommentCount(ctx, postIds, getCommentCountCh)

	authorPreviews, err := mgr.playerService.GetPlayerPreviews(ctx, userIDs)
	if err != nil {
		return models.Newsfeed{}, err
	}

	getLikeCountRes := <-getLikeCountCh
	if getLikeCountRes.Err != nil {
		return models.Newsfeed{}, getLikeCountRes.Err
	}

	getCommentCountRes := <-getCommentCountCh
	if getCommentCountRes.Err != nil {
		return models.Newsfeed{}, getCommentCountRes.Err
	}

	newCursor := cursor
	for i, post := range posts {
		for _, rawMedia := range post.Medias {
			mediaID, err := uuid.Parse(rawMedia["id"].(string))
			if err != nil {
				return models.Newsfeed{}, errors.WithStack(err)
			}

			url, method, err := mgr.linkGenerator.GenerateDownloadLink(ctx, media.Key(mediaID))

			rawMedia["url"] = url
			rawMedia["method"] = method
		}

		posts[i].Author.Complete(authorPreviews[post.Author.ID])
		posts[i].LikeCount = getLikeCountRes.Data[i]
		posts[i].CommentCount = getCommentCountRes.Data[i]

		newCursor = min(newCursor, post.Timestamp)
	}

	return models.Newsfeed{
		Posts:        posts,
		Cursor:       newCursor,
		HasMorePosts: len(posts) == postPerRequestCount,
	}, nil
}

func (mgr *newsfeedManager) getPostsLikeCount(
	ctx context.Context,
	postIds []uuid.UUID,
	ch chan<- result.Result[[]int64],
) {
	likeCount, err := mgr.postRepository.GetPostsLikeCount(ctx, postIds)
	if err != nil {
		ch <- result.Result[[]int64]{Err: err}
	}

	ch <- result.Result[[]int64]{Data: likeCount}
}

func (mgr *newsfeedManager) getPostsCommentCount(
	ctx context.Context,
	postIds []uuid.UUID,
	ch chan<- result.Result[[]int64],
) {
	commentCount, err := mgr.postRepository.GetPostsCommentCount(ctx, postIds)
	if err != nil {
		ch <- result.Result[[]int64]{Err: err}
	}

	ch <- result.Result[[]int64]{Data: commentCount}
}
