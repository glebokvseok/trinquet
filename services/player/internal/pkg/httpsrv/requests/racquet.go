package requests

import "github.com/move-mates/trinquet/services/player/internal/pkg/models"

type UpdateRacquetProfilesRequest struct {
	Updates []models.RacquetProfileUpdate `json:"updates" validate:"required,min=1,max=2,dive"`
}
