package responses

import (
	"github.com/move-mates/trinquet/services/notification/internal/pkg/models"
)

type GetNotificationsResponse struct {
	Notifications []*models.Notification `json:"notifications"`
}
