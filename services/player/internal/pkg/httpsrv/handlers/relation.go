package handlers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	apierrors "github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"net/http"
	"strconv"
)

type PlayerRelationHandler struct {
	relationManager managers.RelationManager
	validator       *validator.Validate
}

func providePlayerRelationHandler(
	relationManager managers.RelationManager,
	validator *validator.Validate,
) *PlayerRelationHandler {
	return &PlayerRelationHandler{
		relationManager: relationManager,
		validator:       validator,
	}
}

func (handler *PlayerRelationHandler) FollowPlayer(ctx echo.Context) error {
	var request requests.FollowPlayerRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	err = handler.relationManager.FollowUser(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.PlayerId,
		request.FollowedOn,
	)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *PlayerRelationHandler) UnfollowPlayer(ctx echo.Context) error {
	var request requests.UnfollowPlayerRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	err = handler.relationManager.UnfollowUser(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.PlayerId,
	)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *PlayerRelationHandler) GetFollowing(ctx echo.Context) error {
	return handler.getFollows(ctx, handler.relationManager.GetFollowingSection)
}

func (handler *PlayerRelationHandler) GetFollowers(ctx echo.Context) error {
	return handler.getFollows(ctx, handler.relationManager.GetFollowersSection)
}

func (handler *PlayerRelationHandler) getFollows(
	ctx echo.Context,
	getFollowSectionFunc func(context.Context, uuid.UUID, uuid.UUID, models.FollowSort) (models.FollowSection, error),
) error {
	playerID := auth.GetUserID(ctx)
	if playerIDParam := ctx.QueryParam("id"); playerIDParam != "self" {
		var err error
		if playerID, err = uuid.Parse(playerIDParam); err != nil {
			return apierrors.NewInvalidPlayerIDFormatError(playerIDParam)
		}
	}

	cursor, err := strconv.ParseInt(ctx.QueryParam("cursor"), 10, 64)
	if err != nil {
		return apierrors.NewInvalidCursorFormatError()
	}

	sort := models.FollowSort{
		Type:   models.FollowSortType(ctx.QueryParam("sort")),
		Cursor: cursor,
	}

	followSection, err := getFollowSectionFunc(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		playerID,
		sort,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, followSection)
}
