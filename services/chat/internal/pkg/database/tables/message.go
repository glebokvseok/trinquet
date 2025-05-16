package tables

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"time"
)

const (
	MessageTableName = "chat_service.message"
)

type Message struct {
	ID                uuid.UUID            `gorm:"column:id;type:uuid;primaryKey"`
	ChatID            uuid.UUID            `gorm:"column:chat_id;type:uuid"`
	UserID            uuid.UUID            `gorm:"column:user_id;type:uuid"`
	Type              domain.MessageType   `gorm:"type"`
	Status            domain.MessageStatus `gorm:"status"`
	StatusDescription string               `gorm:"status_description"`
	Content           datatypes.JSON       `gorm:"content"`
	ClientCreatedOn   time.Time            `gorm:"column:client_created_on"`
	ClientModifiedOn  time.Time            `gorm:"column:client_modified_on"`
	CreatedOn         time.Time            `gorm:"column:created_on"`
	ModifiedOn        time.Time            `gorm:"column:modified_on"`
}

func (*Message) TableName() string {
	return MessageTableName
}

func (message *Message) ToMessageModel() (*models.Message, error) {
	var content map[string]any
	err := json.Unmarshal(message.Content, &content)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.Message{
		ID:         message.ID,
		UserID:     message.UserID,
		Type:       message.Type,
		Status:     message.Status,
		Content:    content,
		CreatedOn:  message.ClientCreatedOn,
		ModifiedOn: message.ModifiedOn,
	}, nil
}
