package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	apierrors "github.com/move-mates/trinquet/services/newsfeed/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type PostMediaHandler struct {
	mediaManager managers.MediaManager
	validator    *validator.Validate
}

func providePostMediaHandler(
	mediaManager managers.MediaManager,
	validator *validator.Validate,
) *PostMediaHandler {
	return &PostMediaHandler{
		mediaManager: mediaManager,
		validator:    validator,
	}
}

func (handler *PostMediaHandler) UploadPostMedia(ctx echo.Context) error {
	rawMedia, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return errors.WithStack(err)
	}

	var mediaJson map[string]any
	err = json.Unmarshal(rawMedia, &mediaJson)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	mediaType, ok := mediaJson["media_type"].(models.MediaType)
	if !ok {
		return apierrors.NewInvalidMediaTypeError()
	}

	var media any
	switch mediaType {
	case models.PhotoMediaType:
		media, err = parseMedia[models.Photo](rawMedia, handler.validator)
	case models.VideoMediaType:
		media, err = parseMedia[models.Video](rawMedia, handler.validator)
	default:
		return apierrors.NewUnsupportedMediaTypeError(mediaType)
	}

	if err != nil {
		return err
	}

	obj, err := handler.mediaManager.UploadMedia(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		media,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, obj)
}

func parseMedia[TMedia any](rawMedia []byte, validator *validator.Validate) (TMedia, error) {
	var media TMedia
	err := json.Unmarshal(rawMedia, &media)
	if err != nil {
		return media, cmnerrors.NewRequestBodyParsingError(err)
	}

	err = validator.Struct(media)
	if err != nil {
		return media, cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	return media, nil
}
