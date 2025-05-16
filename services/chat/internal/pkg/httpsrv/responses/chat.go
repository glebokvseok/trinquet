package responses

import (
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
)

type GetChatMessagesResponse struct {
	Messages []*models.Message `json:"messages"`
}
