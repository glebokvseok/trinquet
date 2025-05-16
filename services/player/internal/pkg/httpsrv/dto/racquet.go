package dto

import (
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
)

type RacquetProfile struct {
	SportType  domain.RacquetSportType `json:"sport_type"`
	BestHand   domain.BestHand         `json:"best_hand"`
	CourtSide  domain.CourtSide        `json:"court_side"`
	Rating     int                     `json:"rating"`
	MatchCount int                     `json:"match_count"`
	WinCount   int                     `json:"win_count"`
	LossCount  int                     `json:"lost_count"`
}
