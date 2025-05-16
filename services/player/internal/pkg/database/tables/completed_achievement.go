package tables

import (
	"github.com/google/uuid"
	"time"
)

const (
	CompletedAchievementTableName = "player_service.completed_achievement"
)

type CompletedAchievement struct {
	playerID      uuid.UUID `gorm:"column:player_id;type:uuid"`
	achievementID uuid.UUID `gorm:"column:achievement_id;type:uuid"`
	CompletedOn   time.Time `gorm:"column:completed_on"`
}

func (*CompletedAchievement) TableName() string {
	return CompletedAchievementTableName
}
