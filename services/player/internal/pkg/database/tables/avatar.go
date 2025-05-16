package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"time"
)

const (
	AvatarTableName = "player_service.avatar"
)

type Avatar struct {
	AvatarID   uuid.UUID `gorm:"column:avatar_id;type:uuid;primaryKey"`
	PlayerID   uuid.UUID `gorm:"column:player_id;type:uuid"`
	MimeType   string    `gorm:"column:mime_type"`
	IsUsed     bool      `gorm:"column:is_used"`
	CreatedOn  time.Time `gorm:"column:created_on"`
	ModifiedOn time.Time `gorm:"column:modified_on"`
}

func (avatar *Avatar) TableName() string {
	return AvatarTableName
}

func (avatar *Avatar) ToAvatarModel() *models.Avatar {
	return &models.Avatar{
		ID:       avatar.AvatarID,
		MimeType: avatar.MimeType,
	}
}
