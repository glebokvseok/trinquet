package repos

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/events"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, message events.NewMessageEvent) error
	GetChatMessages(ctx context.Context, chatID uuid.UUID, since time.Time) ([]*models.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func ProvideMessageRepository(
	db *gorm.DB,
) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (repo *messageRepository) SaveMessage(
	ctx context.Context,
	userID uuid.UUID,
	chatID uuid.UUID,
	message events.NewMessageEvent,
) error {
	content, err := json.Marshal(message.Content)
	if err != nil {
		return errors.WithStack(err)
	}

	currentTime := time.Now()
	err = repo.db.
		WithContext(ctx).
		Create(&tables.Message{
			ID:               uuid.New(),
			ChatID:           chatID,
			UserID:           userID,
			Type:             domain.StandardMessage,
			Status:           domain.SentMessageStatus,
			Content:          content,
			ClientCreatedOn:  message.SentOn,
			ClientModifiedOn: message.SentOn,
			CreatedOn:        currentTime,
			ModifiedOn:       currentTime,
		}).
		Error

	return errors.WithStack(err)
}

func (repo *messageRepository) GetChatMessages(
	ctx context.Context,
	chatID uuid.UUID,
	since time.Time,
) ([]*models.Message, error) {
	var messages []*tables.Message
	if err := repo.db.
		WithContext(ctx).
		Where("chat_id = ? and client_created_on > ?", chatID, since).
		Order("client_created_on asc").
		Find(&messages).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return slice.MapErr(
		messages,
		func(message *tables.Message) (*models.Message, error) {
			return message.ToMessageModel()
		},
	)
}
