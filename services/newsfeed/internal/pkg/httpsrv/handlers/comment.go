package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	apierrors "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/managers"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	commentManager managers.CommentManager
}

func provideCommentHandler(
	commentManager managers.CommentManager,
) *CommentHandler {
	return &CommentHandler{
		commentManager: commentManager,
	}
}

func (handler *CommentHandler) GetCommentSection(ctx echo.Context) error {
	postID, err := uuid.Parse(ctx.QueryParam("id"))
	if err != nil {
		return apierrors.NewInvalidPostIDFormatError()
	}

	cursor, err := strconv.ParseInt(ctx.QueryParam("cursor"), 10, 64)
	if err != nil {
		return apierrors.NewInvalidCursorFormatError()
	}

	commentSection, err := handler.commentManager.GetCommentSection(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		postID,
		cursor,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, commentSection)
}
