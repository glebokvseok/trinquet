package models

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
)

type ExternalChat struct {
	ID   uuid.UUID       `json:"id" validate:"required"`
	Type domain.ChatType `json:"type" validate:"required,min=1,max=2"`
}

type InternalChat struct {
	ID   uuid.UUID
	Type domain.ChatType
	Key  string
}
