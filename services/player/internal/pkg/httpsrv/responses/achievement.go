package responses

import "github.com/move-mates/trinquet/services/player/internal/pkg/models"

type GetAchievementsResponse struct {
	Achievements []*models.Achievement `json:"achievements"`
}
