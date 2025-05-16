package managers

import (
	"context"
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/court/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/court/internal/pkg/helpers/media"
	"github.com/move-mates/trinquet/services/court/internal/pkg/models"
)

type CourtManager interface {
	SearchCourts(ctx context.Context, filter models.SearchFilter) (models.SearchResult, error)
	GetCourtMedia(ctx context.Context, courtID uuid.UUID) (models.CourtMedia, error)
}

type courtManager struct {
	courtRepository repos.CourtRepository
	linkGenerator   s3.LinkGenerator
}

func ProvideCourtManager(
	courtRepository repos.CourtRepository,
	linkGenerator s3.LinkGenerator,
) CourtManager {
	return &courtManager{
		courtRepository: courtRepository,
		linkGenerator:   linkGenerator,
	}
}

func (mgr *courtManager) SearchCourts(
	ctx context.Context,
	filter models.SearchFilter,
) (models.SearchResult, error) {
	courts, err := mgr.courtRepository.SearchCourt(ctx, filter)
	if err != nil {
		return models.SearchResult{}, err
	}

	for i, court := range courts {
		courts[i].Preview.URL, courts[i].Preview.Method, err =
			mgr.linkGenerator.GenerateDownloadLink(ctx, media.Key(court.ID, court.Preview.ID))

		if err != nil {
			return models.SearchResult{}, err
		}
	}

	return models.SearchResult{
		Courts: courts,
	}, nil
}

func (mgr *courtManager) GetCourtMedia(
	ctx context.Context,
	courtID uuid.UUID,
) (models.CourtMedia, error) {
	medias, err := mgr.courtRepository.GetCourtMedia(ctx, courtID)
	if err != nil {
		return models.CourtMedia{}, err
	}

	for i := range medias {
		medias[i].URL, medias[i].Method, err =
			mgr.linkGenerator.GenerateDownloadLink(ctx, media.Key(courtID, medias[i].ID))

		if err != nil {
			return models.CourtMedia{}, err
		}
	}

	return models.CourtMedia{
		Medias: medias,
	}, nil
}
