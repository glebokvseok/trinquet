package models

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/dto"
	"time"
)

type RacquetProfile struct {
	ID         uuid.UUID
	SportType  domain.RacquetSportType
	BestHand   domain.BestHand
	CourtSide  domain.CourtSide
	Rating     int
	MatchCount int
	WinCount   int
	LossCount  int
}

type RacquetProfileUpdate struct {
	SportType domain.RacquetSportType `json:"sport_type" validate:"omitempty,min=1,max=2"`
	BestHand  domain.BestHand         `json:"best_hand" validate:"omitempty,min=1,max=3"`
	CourtSide domain.CourtSide        `json:"court_side" validate:"omitempty,min=1,max=3"`
}

type RacquetMatchUpdate struct {
	SportType     domain.RacquetSportType `json:"sport_type" validate:"required,min=1,max=2"`
	MatchType     domain.RacquetMatchType `json:"match_type" validate:"required,min=1,max=2"`
	IsCompetitive bool                    `json:"is_competitive" validate:"omitempty"`
	CourtID       *uuid.UUID              `json:"court_id" validate:"omitempty"`
	ScheduledOn   time.Time               `json:"scheduled_on" validate:"required"`
}

func (profile *RacquetProfile) ToRacquetProfileDTO() *dto.RacquetProfile {
	return &dto.RacquetProfile{
		SportType:  profile.SportType,
		BestHand:   profile.BestHand,
		CourtSide:  profile.CourtSide,
		Rating:     profile.Rating,
		MatchCount: profile.MatchCount,
		WinCount:   profile.WinCount,
	}
}
