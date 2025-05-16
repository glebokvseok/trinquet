package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"time"
)

type ChatManager interface {
	GetChatMessages(ctx context.Context, userID uuid.UUID, chat models.ExternalChat, since time.Time) ([]*models.Message, error)
}

type chatManager struct {
	chatRepository    repos.ChatRepository
	messageRepository repos.MessageRepository
}

func ProvideChatManager(
	chatRepository repos.ChatRepository,
	messageRepository repos.MessageRepository,
) ChatManager {
	return &chatManager{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
	}
}

func (mgr *chatManager) GetChatMessages(
	ctx context.Context,
	userID uuid.UUID,
	externalChat models.ExternalChat,
	since time.Time,
) ([]*models.Message, error) {
	chat, err := mgr.chatRepository.GetInternalChat(ctx, userID, externalChat)
	if err != nil {
		return nil, err
	}

	return mgr.messageRepository.GetChatMessages(ctx, chat.ID, since)
}
