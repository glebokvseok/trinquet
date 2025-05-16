package managers

import (
	"context"
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/helpers/avatar"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
)

type AvatarManager interface {
	DownloadSelf(ctx context.Context, selfID uuid.UUID) (*models.Avatar, error)
	DownloadOther(ctx context.Context, selfID uuid.UUID, playerID uuid.UUID) (*models.Avatar, error)
	Upload(ctx context.Context, playerID uuid.UUID, mimeType string) (models.Avatar, error)
}

func ProvideAvatarManager(
	linkGenerator s3.LinkGenerator,
	avatarRepository repos.AvatarRepository,
) AvatarManager {
	return &avatarManager{
		linkGenerator:    linkGenerator,
		avatarRepository: avatarRepository,
	}
}

type avatarManager struct {
	linkGenerator    s3.LinkGenerator
	avatarRepository repos.AvatarRepository
}

func (mgr *avatarManager) DownloadSelf(ctx context.Context, selfID uuid.UUID) (*models.Avatar, error) {
	return mgr.download(ctx, selfID, selfID) // TODO: переделать хак для получения собственной аватарки
}

func (mgr *avatarManager) DownloadOther(ctx context.Context, selfID uuid.UUID, playerID uuid.UUID) (*models.Avatar, error) {
	return mgr.download(ctx, selfID, playerID)
}

func (mgr *avatarManager) download(ctx context.Context, selfID uuid.UUID, playerID uuid.UUID) (*models.Avatar, error) {
	av, err := mgr.avatarRepository.GetAvatarByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if av == nil {
		return nil, nil
	}

	url, method, err := mgr.linkGenerator.GenerateDownloadLink(ctx, avatar.Key(playerID, av.ID))
	if err != nil {
		return nil, err
	}

	av.URL, av.Method = url, method

	return av, err
}

func (mgr *avatarManager) Upload(ctx context.Context, playerID uuid.UUID, mimeType string) (models.Avatar, error) {
	avatarID := uuid.New()

	url, method, err := mgr.linkGenerator.GenerateUploadLink(ctx, avatar.Key(playerID, avatarID))
	if err != nil {
		return models.Avatar{}, err
	}

	if err = mgr.avatarRepository.AddAvatar(ctx, playerID, avatarID, mimeType); err != nil {
		return models.Avatar{}, err
	}

	return models.Avatar{
		ID:     avatarID,
		URL:    url,
		Method: method,
	}, nil
}
