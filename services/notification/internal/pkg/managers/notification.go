package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/models"
)

type NotificationManager interface {
	SaveNotification(ctx context.Context, userID uuid.UUID, notification models.Notification) error
	GetNotifications(ctx context.Context, userID uuid.UUID, timestamp int64) ([]*models.Notification, error)
}

type notificationManager struct {
	notificationRepository repos.NotificationRepository
}

func ProvideNotificationManager(
	notificationRepository repos.NotificationRepository,
) NotificationManager {
	return &notificationManager{
		notificationRepository: notificationRepository,
	}
}

func (mgr *notificationManager) SaveNotification(
	ctx context.Context,
	userID uuid.UUID,
	notification models.Notification,
) error {
	return mgr.notificationRepository.SaveNotification(ctx, userID, notification)
}

func (mgr *notificationManager) GetNotifications(
	ctx context.Context,
	userID uuid.UUID,
	timestamp int64,
) ([]*models.Notification, error) {
	return mgr.notificationRepository.GetNotifications(ctx, userID, timestamp)
}
