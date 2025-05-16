package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	apierrors "github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type RacquetProfileRepository interface {
	GetProfile(ctx context.Context, playerID uuid.UUID, sportType domain.RacquetSportType) (*models.RacquetProfile, error)
	GetProfiles(ctx context.Context, playerID uuid.UUID) ([]*models.RacquetProfile, error)
	CreateProfilesTx(ctx context.Context, playerID uuid.UUID) error
	UpdateProfiles(ctx context.Context, playerID uuid.UUID, updates []models.RacquetProfileUpdate) error
}

type racquetProfileRepository struct {
	db *gorm.DB
}

func ProvideRacquetProfileRepository(
	db *gorm.DB,
) RacquetProfileRepository {
	return &racquetProfileRepository{
		db: db,
	}
}

func (repo *racquetProfileRepository) GetProfile(
	ctx context.Context,
	playerID uuid.UUID,
	sportType domain.RacquetSportType,
) (*models.RacquetProfile, error) {
	var profile = new(tables.RacquetProfile)
	if err := repo.db.
		WithContext(ctx).
		Where("player_id = ? and sport_type = ?", playerID, sportType).
		Take(profile).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewRacquetProfileNotFoundError(playerID, sportType)
		}

		return nil, errors.WithStack(err)
	}

	return profile.ToRacquetProfileModel(), nil
}

func (repo *racquetProfileRepository) GetProfiles(
	ctx context.Context,
	playerID uuid.UUID,
) ([]*models.RacquetProfile, error) {
	var profiles []*tables.RacquetProfile
	if err := repo.db.
		WithContext(ctx).
		Where("player_id = ?", playerID).
		Order("sport_type asc").
		Find(&profiles).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return slice.Map(
		profiles,
		func(profile *tables.RacquetProfile) *models.RacquetProfile {
			return profile.ToRacquetProfileModel()
		},
	), nil
}

func (repo *racquetProfileRepository) CreateProfilesTx(
	ctx context.Context,
	playerID uuid.UUID,
) error {
	tx, err := getTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	currentTime := time.Now()

	profiles := []tables.RacquetProfile{
		createDefaultRacquetProfile(playerID, domain.Tennis, currentTime),
		createDefaultRacquetProfile(playerID, domain.Padel, currentTime),
	}

	if err = tx.
		WithContext(ctx).
		Create(&profiles).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *racquetProfileRepository) UpdateProfiles(
	ctx context.Context,
	playerID uuid.UUID,
	updates []models.RacquetProfileUpdate,
) (returnErr error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return errors.WithStack(tx.Error)
	}

	defer psql.TransactionFinalizer(tx, &returnErr)

	currentTime := time.Now()

	for _, update := range updates {
		res := tx.
			WithContext(ctx).
			Table(tables.RacquetProfileTableName).
			Where("player_id = ? and sport_type = ?", playerID, update.SportType).
			Updates(map[string]any{
				"best_hand":   update.BestHand,
				"court_side":  update.CourtSide,
				"modified_on": currentTime,
			})

		if res.Error != nil {
			return errors.WithStack(res.Error)
		}

		if res.RowsAffected != 1 {
			return apierrors.NewRacquetProfileNotFoundError(playerID, update.SportType)
		}
	}

	return errors.WithStack(tx.Commit().Error)
}

func createDefaultRacquetProfile(
	playerID uuid.UUID,
	sportType domain.RacquetSportType,
	currentTime time.Time,
) tables.RacquetProfile {
	return tables.RacquetProfile{
		ID:         uuid.New(),
		PlayerID:   playerID,
		SportType:  sportType,
		Rating:     domain.DefaultRacquetRating,
		CreatedOn:  currentTime,
		ModifiedOn: currentTime,
	}
}
