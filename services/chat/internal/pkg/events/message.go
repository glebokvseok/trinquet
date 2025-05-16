package events

import (
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"time"
)

type MessageEventType = string

const (
	NewMessage MessageEventType = "new_message"
)

type NewMessageEvent struct {
	Chat    models.ExternalChat   `json:"chat" validate:"required"`
	Content models.MessageContent `json:"content" validate:"required"`
	SentOn  time.Time             `json:"sent_on" validate:"required"`
}
