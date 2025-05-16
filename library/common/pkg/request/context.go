package request

import (
	"context"
	"github.com/google/uuid"
)

type contextKey string

const (
	requestIDContextKey contextKey = "request_id"
)

func CreateContext(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIDContextKey, id)
}

func GetIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(requestIDContextKey).(uuid.UUID)
}

func TryGetIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(requestIDContextKey).(uuid.UUID)

	return id, ok
}
