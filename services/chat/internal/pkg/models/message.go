package models

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
	"time"
)

type Message struct {
	ID         uuid.UUID            `json:"id"`
	UserID     uuid.UUID            `json:"user_id"`
	Type       domain.MessageType   `json:"type"`
	Status     domain.MessageStatus `json:"status"`
	Content    map[string]any       `json:"content"`
	CreatedOn  time.Time            `json:"created_on"`
	ModifiedOn time.Time            `json:"modified_on"`
}

type MessageContent struct {
	Text string `json:"text" validate:"required"`
}
