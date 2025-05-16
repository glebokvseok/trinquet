package tables

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

const (
	NotificationTableName = "notification_service.user_notification"
)

type Notification struct {
	ID        uuid.UUID               `gorm:"column:notification_id;type:uuid"`
	UserID    uuid.UUID               `gorm:"column:user_id;type:uuid;primaryKey"`
	Type      domain.NotificationType `gorm:"column:notification_type;type:uuid;primaryKey"`
	Data      datatypes.JSON          `gorm:"column:data;type:jsonb"`
	Timestamp int64                   `gorm:"column:timestamp"`
}

func (*Notification) TableName() string {
	return NotificationTableName
}

func (notification *Notification) ToNotificationModel() (*models.Notification, error) {
	var data map[string]any
	err := json.Unmarshal(notification.Data, &data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.Notification{
		Type:      notification.Type,
		Data:      data,
		Timestamp: notification.Timestamp,
	}, nil
}
