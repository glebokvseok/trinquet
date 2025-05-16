package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"net/http"
)

type RacquetProfileHandler struct {
	profileManager managers.RacquetProfileManager
	validator      *validator.Validate
}

func provideRacquetProfileHandler(
	profileManager managers.RacquetProfileManager,
	validator *validator.Validate,
) *RacquetProfileHandler {
	return &RacquetProfileHandler{
		profileManager: profileManager,
		validator:      validator,
	}
}

func (handler *RacquetProfileHandler) UpdateRacquetProfiles(ctx echo.Context) error {
	var request requests.UpdateRacquetProfilesRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	err = handler.profileManager.UpdateProfiles(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.Updates,
	)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
