package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"net/http"
)

type RacquetMatchHandler struct {
	matchManager managers.RacquetMatchManager
	validator    *validator.Validate
}

func provideRacquetMatchHandler(
	matchManager managers.RacquetMatchManager,
	validator *validator.Validate,
) *RacquetMatchHandler {
	return &RacquetMatchHandler{
		matchManager: matchManager,
		validator:    validator,
	}
}

func (handler *RacquetMatchHandler) CreateMatch(ctx echo.Context) error {
	var match models.RacquetMatchUpdate
	err := ctx.Bind(&match)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(match)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	matchID, err := handler.matchManager.CreateMatch(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		match,
	)

	return ctx.JSON(
		http.StatusOK,
		responses.CreateMatchResponse{
			MatchID: matchID,
		},
	)
}
