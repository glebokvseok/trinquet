package managers

import (
	"context"
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/helpers/avatar"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"time"
)

const (
	followsPerRequestCount = 20
)

type RelationManager interface {
	FollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, followedOn time.Time) error
	UnfollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) error
	GetAllFollowerIDs(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetFollowersSection(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, sort models.FollowSort) (models.FollowSection, error)
	GetFollowingSection(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, sort models.FollowSort) (models.FollowSection, error)
}

func ProvideRelationManager(
	relationRepository repos.RelationRepository,
	playerRepository repos.PlayerRepository,
	linkGenerator s3.LinkGenerator,
) RelationManager {
	return &relationManager{
		relationRepository: relationRepository,
		playerRepository:   playerRepository,
		linkGenerator:      linkGenerator,
	}
}

type relationManager struct {
	relationRepository repos.RelationRepository
	playerRepository   repos.PlayerRepository
	linkGenerator      s3.LinkGenerator
}

func (mgr *relationManager) FollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, followedOn time.Time) error {
	return mgr.relationRepository.FollowUser(ctx, selfID, userID, followedOn)
}

func (mgr *relationManager) UnfollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) error {
	return mgr.relationRepository.UnfollowUser(ctx, selfID, userID)
}

func (mgr *relationManager) GetAllFollowerIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return mgr.relationRepository.GetAllFollowerIDs(ctx, userID)
}

func (mgr *relationManager) GetFollowersSection(
	ctx context.Context,
	selfID uuid.UUID,
	userID uuid.UUID,
	sort models.FollowSort,
) (models.FollowSection, error) {
	getFollowers := func(ctx context.Context) (models.FollowSortResult, error) {
		return mgr.relationRepository.GetFollowers(ctx, userID, sort, followsPerRequestCount)
	}

	return mgr.getFollowSection(ctx, getFollowers)
}

func (mgr *relationManager) GetFollowingSection(
	ctx context.Context,
	selfID uuid.UUID,
	userID uuid.UUID,
	sort models.FollowSort,
) (models.FollowSection, error) {
	getFollowing := func(ctx context.Context) (models.FollowSortResult, error) {
		return mgr.relationRepository.GetFollowing(ctx, userID, sort, followsPerRequestCount)
	}

	return mgr.getFollowSection(ctx, getFollowing)
}

func (mgr *relationManager) getFollowSection(
	ctx context.Context,
	getFollowSortResultFunc func(ctx context.Context) (models.FollowSortResult, error),
) (models.FollowSection, error) {
	res, err := getFollowSortResultFunc(ctx)
	if err != nil {
		return models.FollowSection{}, err
	}

	ids := make([]uuid.UUID, len(res.Follows))
	for i, follow := range res.Follows {
		ids[i] = follow.ID
	}

	basePreviews, err := mgr.playerRepository.GetBasePlayerPreviews(ctx, ids)
	if err != nil {
		return models.FollowSection{}, err
	}

	previews := make([]models.FollowPreview, len(basePreviews))
	for i, basePreview := range basePreviews {
		previews[i].BasePlayerPreview = basePreview
		previews[i].Following = res.Follows[i].Following
		previews[i].FollowingBack = res.Follows[i].FollowingBack
		if basePreview.Avatar != nil {
			previews[i].Avatar.URL, previews[i].Avatar.Method, err =
				mgr.linkGenerator.GenerateDownloadLink(ctx, avatar.Key(basePreview.PlayerID, basePreview.Avatar.ID))

			if err != nil {
				return models.FollowSection{}, err
			}
		}
	}

	return models.FollowSection{
		Follows:        previews,
		Cursor:         res.Cursor,
		HasMoreFollows: len(previews) == followsPerRequestCount,
	}, nil
}
