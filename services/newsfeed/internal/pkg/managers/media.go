package managers

import (
	"context"
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/helpers/media"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
)

type MediaManager interface {
	UploadMedia(ctx context.Context, userID uuid.UUID, mediaData any) (models.Media, error)
}

func ProvideMediaManager(
	linkGenerator s3.LinkGenerator,
	mediaRepository repos.MediaRepository,
) MediaManager {
	return &mediaManager{
		linkGenerator:   linkGenerator,
		mediaRepository: mediaRepository,
	}
}

type mediaManager struct {
	linkGenerator   s3.LinkGenerator
	mediaRepository repos.MediaRepository
}

func (mgr *mediaManager) UploadMedia(ctx context.Context, userID uuid.UUID, mediaData any) (models.Media, error) {
	mediaID := uuid.New()

	err := mgr.mediaRepository.AddMedia(ctx, userID, mediaID, mediaData)
	if err != nil {
		return models.Media{}, err
	}

	url, method, err := mgr.linkGenerator.GenerateUploadLink(ctx, media.Key(mediaID))
	if err != nil {
		return models.Media{}, err
	}

	return models.Media{
		ID:     mediaID,
		URL:    url,
		Method: method,
	}, nil
}
