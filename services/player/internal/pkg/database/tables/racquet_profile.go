package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"time"
)

const (
	RacquetProfileTableName = "player_service.racquet_profile"
)

type RacquetProfile struct {
	ID         uuid.UUID               `gorm:"column:id;type:uuid;primaryKey"`
	PlayerID   uuid.UUID               `gorm:"column:player_id;type:uuid"`
	SportType  domain.RacquetSportType `gorm:"column:sport_type"`
	BestHand   domain.BestHand         `gorm:"column:best_hand"`
	CourtSide  domain.CourtSide        `gorm:"column:court_side"`
	Rating     int                     `gorm:"column:rating"`
	MatchCount int                     `gorm:"column:match_count"`
	WinCount   int                     `gorm:"column:match_count"`
	LossCount  int                     `gorm:"column:match_count"`
	CreatedOn  time.Time               `gorm:"column:created_on"`
	ModifiedOn time.Time               `gorm:"column:modified_on"`
}

func (*RacquetProfile) TableName() string {
	return RacquetProfileTableName
}

func (profile *RacquetProfile) ToRacquetProfileModel() *models.RacquetProfile {
	return &models.RacquetProfile{
		ID:         profile.ID,
		SportType:  profile.SportType,
		BestHand:   profile.BestHand,
		CourtSide:  profile.CourtSide,
		Rating:     profile.Rating,
		MatchCount: profile.MatchCount,
		WinCount:   profile.WinCount,
		LossCount:  profile.LossCount,
	}
}
