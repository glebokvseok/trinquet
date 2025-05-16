package repos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type ChatRepository interface {
	GetInternalChat(ctx context.Context, userID uuid.UUID, chat models.ExternalChat) (*models.InternalChat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func ProvideChatRepository(
	db *gorm.DB,
) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (repo *chatRepository) GetInternalChat(
	ctx context.Context,
	userID uuid.UUID,
	externalChat models.ExternalChat,
) (*models.InternalChat, error) {
	key, err := getChatKey(userID, externalChat.ID, externalChat.Type)
	if err != nil {
		return nil, err
	}

	chat := new(tables.Chat)
	if err = repo.db.
		WithContext(ctx).
		Where("key = ?", key).
		First(chat).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repo.createChat(ctx, key, externalChat.Type)
		}
	}

	return chat.ToInternalChatModel(), nil
}

func (repo *chatRepository) createChat(
	ctx context.Context,
	key string,
	chatType domain.ChatType,
) (*models.InternalChat, error) {
	chat := tables.Chat{
		ID:        uuid.New(),
		Key:       key,
		Type:      chatType,
		CreatedOn: time.Now(),
	}

	if err := repo.db.
		WithContext(ctx).
		Create(&chat).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return chat.ToInternalChatModel(), nil
}

func getChatKey(
	userID uuid.UUID,
	externalID uuid.UUID,
	chatType domain.ChatType,
) (string, error) {
	switch chatType {
	case domain.PersonalChat:
		return fmt.Sprintf(
			"%s_%s",
			min(userID.String(), externalID.String()),
			max(userID.String(), externalID.String()),
		), nil
	default:
		return "", errors.Errorf("unknown chat type %d", chatType)
	}
}
