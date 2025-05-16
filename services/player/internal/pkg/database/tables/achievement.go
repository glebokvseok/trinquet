package tables

import (
	"github.com/google/uuid"
	"time"
)

const (
	AchievementTableName = "player_service.achievement"
)

type Achievement struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	CreatedOn   time.Time `gorm:"column:created_on"`
	ModifiedOn  time.Time `gorm:"column:modified_on"`
}

func (*Achievement) TableName() string {
	return AchievementTableName
}
