package managers

import (
	"context"
	"github.com/google/uuid"
	txmgr "github.com/move-mates/trinquet/library/database/pkg/psql/managers"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
)

type RacquetMatchManager interface {
	CreateMatch(ctx context.Context, playerID uuid.UUID, match models.RacquetMatchUpdate) (uuid.UUID, error)
}

type racquetMatchManager struct {
	transactionManager txmgr.TransactionManager
	matchRepository    repos.RacquetMatchRepository
	profileRepository  repos.RacquetProfileRepository
}

func ProvideRacquetMatchManager(
	transactionManager txmgr.TransactionManager,
	matchRepository repos.RacquetMatchRepository,
	profileRepository repos.RacquetProfileRepository,
) RacquetMatchManager {
	return &racquetMatchManager{
		transactionManager: transactionManager,
		matchRepository:    matchRepository,
		profileRepository:  profileRepository,
	}
}

func (mgr *racquetMatchManager) CreateMatch(
	ctx context.Context,
	playerID uuid.UUID,
	match models.RacquetMatchUpdate,
) (uuid.UUID, error) {
	profile, err := mgr.profileRepository.GetProfile(ctx, playerID, match.SportType)
	if err != nil {
		return uuid.Nil, err
	}

	return mgr.matchRepository.CreateMatch(ctx, profile.ID, match)
}
