package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type RacquetMatchRepository interface {
	CreateMatch(ctx context.Context, profileID uuid.UUID, match models.RacquetMatchUpdate) (uuid.UUID, error)
}

type racquetMatchRepository struct {
	db *gorm.DB
}

func ProvideRacquetMatchRepository(
	db *gorm.DB,
) RacquetMatchRepository {
	return &racquetMatchRepository{
		db: db,
	}
}

func (repo *racquetMatchRepository) CreateMatch(
	ctx context.Context,
	profileID uuid.UUID,
	match models.RacquetMatchUpdate,
) (uuid.UUID, error) {
	matchID, currentTime := uuid.New(), time.Now()

	if err := repo.db.
		WithContext(ctx).
		Create(&tables.RacquetMatch{
			ID:            matchID,
			OwnerID:       profileID,
			SportType:     match.SportType,
			MatchType:     match.MatchType,
			IsCompetitive: match.IsCompetitive,
			CourtID:       match.CourtID,
			ScheduledOn:   match.ScheduledOn,
			State:         domain.RacquetMatchCreated,
			CreatedOn:     currentTime,
			ModifiedOn:    currentTime,
		}).
		Error; err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	return matchID, nil
}
