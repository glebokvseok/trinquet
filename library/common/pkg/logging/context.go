package logging

import (
	"github.com/move-mates/trinquet/library/common/pkg/request"
	"github.com/sirupsen/logrus"
)

const (
	requestIDKey string = "request_id"
)

type ContextHook struct{}

func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		if requestID, ok := request.TryGetIDFromContext(entry.Context); ok {
			entry.Data[requestIDKey] = requestID.String()
		}
	}

	return nil
}
