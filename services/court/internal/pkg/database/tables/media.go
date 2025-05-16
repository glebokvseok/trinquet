package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/court/internal/pkg/models"
	"time"
)

const (
	CourtMediaTableName = "court_service.court_media"
)

type Media struct {
	CourtID   uuid.UUID `gorm:"column:court_id;type:uuid"`
	MediaID   uuid.UUID `gorm:"column:media_id;type:uuid"`
	IsPreview bool      `gorm:"column:is_preview"`
	MimeType  string    `gorm:"column:mime_type"`
	CreatedOn time.Time `gorm:"column:created_on"`
}

func (*Media) TableName() string {
	return CourtMediaTableName
}

func (media *Media) ToMediaModel() *models.Media {
	return &models.Media{
		ID:       media.MediaID,
		MimeType: media.MimeType,
	}
}
