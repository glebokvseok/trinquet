package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	apierrors "github.com/move-mates/trinquet/services/court/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/court/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/court/internal/pkg/models"
	"net/http"
)

type CourtHandler struct {
	courtManager managers.CourtManager
	validator    *validator.Validate
}

func provideCourtHandler(
	courtManager managers.CourtManager,
	validator *validator.Validate,
) *CourtHandler {
	return &CourtHandler{
		courtManager: courtManager,
		validator:    validator,
	}
}

func (handler *CourtHandler) SearchCourts(ctx echo.Context) error {
	var filter models.SearchFilter
	err := ctx.Bind(&filter)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(filter)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	result, err := handler.courtManager.SearchCourts(
		ctx.Request().Context(),
		filter,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (handler *CourtHandler) GetCourtMedia(ctx echo.Context) error {
	courtID, err := uuid.Parse(ctx.QueryParam("id"))
	if err != nil {
		return apierrors.NewInvalidCourtIDFormatError()
	}

	courtMedia, err := handler.courtManager.GetCourtMedia(
		ctx.Request().Context(),
		courtID,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, courtMedia)
}
