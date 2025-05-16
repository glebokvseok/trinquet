package event

import (
	"github.com/google/uuid"
	"time"
)

type Type = string

type Wrapper struct {
	EventType Type      `json:"event_type"`
	RawEvent  []byte    `json:"raw_event"`
	UserID    uuid.UUID `json:"user_id"`
	RequestID uuid.UUID `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}
