package handlers

import (
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	apierrors "github.com/move-mates/trinquet/services/notification/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/notification/internal/pkg/managers"
	"net/http"
	"strconv"
)

type NotificationHandler struct {
	notificationManager managers.NotificationManager
}

func provideNotificationHandler(
	notificationManager managers.NotificationManager,
) *NotificationHandler {
	return &NotificationHandler{
		notificationManager: notificationManager,
	}
}

func (handler *NotificationHandler) GetNotifications(ctx echo.Context) error {
	timestamp, err := strconv.ParseInt(ctx.QueryParam("timestamp"), 10, 64)
	if err != nil {
		return apierrors.NewInvalidTimestampFormatError()
	}

	notifications, err := handler.notificationManager.GetNotifications(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		timestamp,
	)

	return ctx.JSON(
		http.StatusOK,
		responses.GetNotificationsResponse{
			Notifications: notifications,
		},
	)
}
