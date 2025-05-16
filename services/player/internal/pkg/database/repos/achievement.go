package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const (
	getAchievementsQuery = `
		SELECT * FROM player_service.achievement a
		LEFT JOIN player_service.completed_achievement ca
		ON a.id = ca.achievement_id
		WHERE ca.player_id = ?;
	`
)

type AchievementRepository interface {
	GetAchievements(ctx context.Context, playerID uuid.UUID) ([]*models.Achievement, error)
}

type achievementRepository struct {
	db *gorm.DB
}

func ProvideAchievementRepository(
	db *gorm.DB,
) AchievementRepository {
	return &achievementRepository{
		db: db,
	}
}

func (repo *achievementRepository) GetAchievements(
	ctx context.Context,
	playerID uuid.UUID,
) ([]*models.Achievement, error) {
	type achievement struct {
		tables.Achievement
		*tables.CompletedAchievement
	}

	var achievements []*achievement
	if err := repo.db.
		WithContext(ctx).
		Raw(getAchievementsQuery, playerID).
		Scan(&achievements).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return slice.Map(
		achievements,
		func(achievement *achievement) *models.Achievement {
			var completedOn *time.Time
			isCompleted := achievement.CompletedAchievement == nil
			if isCompleted {
				completedOn = &achievement.CompletedOn
			}

			return &models.Achievement{
				Name:        achievement.Name,
				Code:        achievement.Code,
				Description: achievement.Description,
				IsCompleted: isCompleted,
				CompletedOn: completedOn,
			}
		},
	), nil
}
