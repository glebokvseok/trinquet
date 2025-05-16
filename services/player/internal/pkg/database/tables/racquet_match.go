package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"time"
)

const (
	RacquetMatchTableName = "player_service.racquet_match"
)

type RacquetMatch struct {
	ID            uuid.UUID                `gorm:"column:id;type:uuid;primaryKey"`
	OwnerID       uuid.UUID                `gorm:"column:owner_id;type:uuid"`
	SportType     domain.RacquetSportType  `gorm:"column:sport_type"`
	MatchType     domain.RacquetMatchType  `gorm:"column:match_type"`
	IsCompetitive bool                     `gorm:"column:is_competitive"`
	CourtID       *uuid.UUID               `gorm:"column:court_id;type:uuid"`
	State         domain.RacquetMatchState `gorm:"column:state"`
	StartedOn     *time.Time               `gorm:"column:started_on"`
	FinishedOn    *time.Time               `gorm:"column:finished_on"`
	ScheduledOn   time.Time                `gorm:"column:scheduled_on"`
	CreatedOn     time.Time                `gorm:"column:created_on"`
	ModifiedOn    time.Time                `gorm:"column:modified_on"`
}

func (*RacquetMatch) TableName() string {
	return RacquetMatchTableName
}

func (match *RacquetMatch) ToRacquetMatchModel() {

}
