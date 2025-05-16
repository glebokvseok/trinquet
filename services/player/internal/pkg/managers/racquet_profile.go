package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
)

type RacquetProfileManager interface {
	UpdateProfiles(ctx context.Context, playerID uuid.UUID, updates []models.RacquetProfileUpdate) error
}

type racquetProfileManager struct {
	profileRepository repos.RacquetProfileRepository
}

func ProvideRacquetProfileManager(
	profileRepository repos.RacquetProfileRepository,
) RacquetProfileManager {
	return &racquetProfileManager{
		profileRepository: profileRepository,
	}
}

func (mgr *racquetProfileManager) UpdateProfiles(
	ctx context.Context,
	playerID uuid.UUID,
	updates []models.RacquetProfileUpdate,
) error {
	return mgr.profileRepository.UpdateProfiles(ctx, playerID, updates)
}
