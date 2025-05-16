package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"net/http"
)

type AvatarHandler struct {
	avatarManager managers.AvatarManager
	validator     *validator.Validate
}

func provideAvatarHandler(
	avatarManager managers.AvatarManager,
	validator *validator.Validate,
) *AvatarHandler {
	return &AvatarHandler{
		avatarManager: avatarManager,
		validator:     validator,
	}
}

func (handler *AvatarHandler) DownloadSelf(ctx echo.Context) error {
	avatar, err := handler.avatarManager.DownloadSelf(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, avatar)
}

func (handler *AvatarHandler) DownloadOther(ctx echo.Context) error {
	playerID, err := uuid.Parse(ctx.Param(PlayerIDParamName))
	if err != nil {
		return errors.NewInvalidPlayerIDFormatError(ctx.Param(PlayerIDParamName))
	}

	avatar, err := handler.avatarManager.DownloadOther(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		playerID,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, avatar)
}

func (handler *AvatarHandler) Upload(ctx echo.Context) error {
	var request requests.UploadAvatarRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	avatar, err := handler.avatarManager.Upload(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.MimeType,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, avatar)
}
