package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
)

type AchievementManager interface {
	GetAchievements(ctx context.Context, playerID uuid.UUID) ([]*models.Achievement, error)
}

type achievementManager struct {
	achievementRepository repos.AchievementRepository
}

func ProvideAchievementManager(
	achievementRepository repos.AchievementRepository,
) AchievementManager {
	return &achievementManager{
		achievementRepository: achievementRepository,
	}
}

func (mgr *achievementManager) GetAchievements(
	ctx context.Context,
	playerID uuid.UUID,
) ([]*models.Achievement, error) {
	return mgr.achievementRepository.GetAchievements(ctx, playerID)
}
