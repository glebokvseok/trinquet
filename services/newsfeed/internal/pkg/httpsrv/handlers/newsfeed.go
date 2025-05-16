package handlers

import (
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	apierrors "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/managers"
	"net/http"
	"strconv"
)

type NewsfeedHandler struct {
	newsfeedManager managers.NewsfeedManager
}

func provideNewsfeedHandler(
	newsfeedManager managers.NewsfeedManager,
) *NewsfeedHandler {
	return &NewsfeedHandler{
		newsfeedManager: newsfeedManager,
	}
}

func (handler *NewsfeedHandler) GetNewsfeed(ctx echo.Context) error {
	cursor, err := strconv.ParseInt(ctx.QueryParam("cursor"), 10, 64)
	if err != nil {
		return apierrors.NewInvalidCursorFormatError()
	}

	newsfeed, err := handler.newsfeedManager.GetNewsfeed(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		cursor,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, newsfeed)
}
