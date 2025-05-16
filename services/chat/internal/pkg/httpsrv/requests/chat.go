package requests

import (
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"time"
)

type GetChatMessagesRequest struct {
	Chat  models.ExternalChat `json:"chat" validate:"required"`
	Since time.Time           `json:"since" validate:"required"`
}
