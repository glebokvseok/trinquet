package repos

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, userID uuid.UUID, notification models.Notification) error
	GetNotifications(ctx context.Context, userID uuid.UUID, timestamp int64) ([]*models.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func ProvideNotificationRepository(
	db *gorm.DB,
) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (repo *notificationRepository) SaveNotification(
	ctx context.Context,
	userID uuid.UUID,
	notification models.Notification,
) error {
	data, err := json.Marshal(notification.Data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = repo.db.
		WithContext(ctx).
		Create(&tables.Notification{
			ID:        uuid.New(),
			UserID:    userID,
			Type:      notification.Type,
			Data:      data,
			Timestamp: notification.Timestamp,
		}).
		Error

	return errors.WithStack(err)
}

func (repo *notificationRepository) GetNotifications(
	ctx context.Context,
	userID uuid.UUID,
	timestamp int64,
) ([]*models.Notification, error) {
	var notifications []*tables.Notification
	if err := repo.db.
		WithContext(ctx).
		Where("user_id = ? and timestamp < ?", userID, timestamp).
		Find(&notifications).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return slice.MapErr(
		notifications,
		func(notification *tables.Notification) (*models.Notification, error) {
			return notification.ToNotificationModel()
		},
	)
}
