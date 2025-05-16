package models

import "github.com/move-mates/trinquet/services/notification/internal/pkg/domain"

type Notification struct {
	Type      domain.NotificationType `json:"type"`
	Data      map[string]any          `json:"data"`
	Timestamp int64                   `json:"timestamp"`
}
