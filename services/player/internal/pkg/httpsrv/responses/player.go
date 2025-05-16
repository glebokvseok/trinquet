package responses

import (
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/dto"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
)

type GetPlayerResponse struct {
	Info            *dto.Player           `json:"info"`
	RacquetProfiles []*dto.RacquetProfile `json:"racquet_profiles"`
}

type SearchPlayerResponse struct {
	Players []models.SearchPreview `json:"players"`
}
