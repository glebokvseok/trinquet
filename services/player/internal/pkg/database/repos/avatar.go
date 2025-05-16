package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type AvatarRepository interface {
	AddAvatar(ctx context.Context, playerID uuid.UUID, avatarID uuid.UUID, mimeType string) error
	GetAvatarByPlayerID(ctx context.Context, playerID uuid.UUID) (*models.Avatar, error)
	RemoveCurrentAvatarTx(ctx context.Context, playerID uuid.UUID, currentTime time.Time) error
	SetAvatarTx(ctx context.Context, playerID uuid.UUID, avatarID uuid.UUID, currentTime time.Time) error
}

type avatarRepository struct {
	db *gorm.DB
}

func ProvideAvatarRepository(db *gorm.DB) AvatarRepository {
	return &avatarRepository{
		db: db,
	}
}

func (repo *avatarRepository) AddAvatar(
	ctx context.Context,
	playerID uuid.UUID,
	avatarID uuid.UUID,
	mimeType string,
) error {
	currentTime := time.Now()
	avatar := tables.Avatar{
		AvatarID:   avatarID,
		PlayerID:   playerID,
		MimeType:   mimeType,
		IsUsed:     false,
		CreatedOn:  currentTime,
		ModifiedOn: currentTime,
	}

	if err := repo.db.
		WithContext(ctx).
		Create(&avatar).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *avatarRepository) GetAvatarByPlayerID(
	ctx context.Context,
	playerID uuid.UUID,
) (*models.Avatar, error) {
	var avatar = new(tables.Avatar)
	if err := repo.db.
		WithContext(ctx).
		Take(avatar, "player_id = ? and is_used = true", playerID).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.WithStack(err)
	}

	return avatar.ToAvatarModel(), nil
}

func (repo *avatarRepository) RemoveCurrentAvatarTx(
	ctx context.Context,
	playerID uuid.UUID,
	currentTime time.Time,
) error {
	tx, err := getTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.
		WithContext(ctx).
		Table(tables.AvatarTableName).
		Where("player_id = ? and is_used = true", playerID).
		Updates(
			map[string]any{
				"is_used":     false,
				"modified_on": currentTime,
			},
		).Error

	return errors.WithStack(err)
}

func (repo *avatarRepository) SetAvatarTx(
	ctx context.Context,
	playerID uuid.UUID,
	avatarID uuid.UUID,
	currentTime time.Time,
) error {
	tx, err := getTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	res := tx.
		WithContext(ctx).
		Table(tables.AvatarTableName).
		Where("player_id = ? and avatar_id = ?", playerID, avatarID).
		Updates(
			map[string]any{
				"is_used":     true,
				"modified_on": currentTime,
			},
		)

	if res.Error != nil {
		return errors.WithStack(res.Error)
	}

	if res.RowsAffected != 1 {
		return errors.New("failed to set avatar")
	}

	return nil
}
