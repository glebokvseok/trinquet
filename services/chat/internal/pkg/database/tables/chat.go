package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"time"
)

const (
	ChatTableName = "chat_service.chat"
)

type Chat struct {
	ID        uuid.UUID       `gorm:"column:id;type:uuid;primaryKey"`
	Key       string          `gorm:"column:key"`
	Type      domain.ChatType `gorm:"column:type"`
	CreatedOn time.Time       `gorm:"column:created_on"`
}

func (*Chat) TableName() string {
	return ChatTableName
}

func (chat *Chat) ToInternalChatModel() *models.InternalChat {
	return &models.InternalChat{
		ID:   chat.ID,
		Key:  chat.Key,
		Type: chat.Type,
	}
}
