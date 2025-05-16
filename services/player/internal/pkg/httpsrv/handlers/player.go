package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	apierrors "github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/dto"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"net/http"
	"unicode/utf8"
)

const (
	PlayerIDParamName = "id"

	searchQueryMinLength = 4
)

type PlayerHandler struct {
	playerManager      managers.PlayerManager
	achievementManager managers.AchievementManager
	validator          *validator.Validate
}

func providePlayerHandler(
	playerManager managers.PlayerManager,
	achievementManager managers.AchievementManager,
	validator *validator.Validate,
) *PlayerHandler {
	return &PlayerHandler{
		playerManager:      playerManager,
		achievementManager: achievementManager,
		validator:          validator,
	}
}

func (handler *PlayerHandler) GetSelf(ctx echo.Context) error {
	player, racquetProfiles, err := handler.playerManager.GetSelf(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.GetPlayerResponse{
			Info: player.ToPlayerDTO(),
			RacquetProfiles: slice.Map(
				racquetProfiles,
				func(profile *models.RacquetProfile) *dto.RacquetProfile {
					return profile.ToRacquetProfileDTO()
				},
			),
		},
	)
}

func (handler *PlayerHandler) GetOther(ctx echo.Context) error {
	playerID, err := uuid.Parse(ctx.Param(PlayerIDParamName))
	if err != nil {
		return apierrors.NewInvalidPlayerIDFormatError(ctx.Param(PlayerIDParamName))
	}

	player, racquetProfiles, err := handler.playerManager.GetOther(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		playerID,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.GetPlayerResponse{
			Info: player.ToPlayerDTO(),
			RacquetProfiles: slice.Map(
				racquetProfiles,
				func(profile *models.RacquetProfile) *dto.RacquetProfile {
					return profile.ToRacquetProfileDTO()
				},
			),
		},
	)
}

func (handler *PlayerHandler) SearchPlayer(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	if utf8.RuneCountInString(query) < searchQueryMinLength {
		return apierrors.NewShortSearchQueryError(searchQueryMinLength)
	}

	players, err := handler.playerManager.SearchPlayers(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		query,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.SearchPlayerResponse{
			Players: players,
		},
	)
}

func (handler *PlayerHandler) CreatePlayer(ctx echo.Context) error {
	var request requests.CreatePlayerRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	err = handler.playerManager.CreatePlayer(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.Username,
	)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *PlayerHandler) UpdatePlayer(ctx echo.Context) error {
	var update models.PlayerUpdate
	err := ctx.Bind(&update)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(update)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	err = handler.playerManager.UpdatePlayer(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		update,
	)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *PlayerHandler) GetAchievements(ctx echo.Context) error {
	achievements, err := handler.achievementManager.GetAchievements(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.GetAchievementsResponse{
			Achievements: achievements,
		},
	)
}
