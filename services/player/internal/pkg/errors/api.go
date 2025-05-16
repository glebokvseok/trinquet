package errors

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"net/http"
)

const (
	InvalidPlayerIDFormatError = "INVALID_PLAYER_ID_FORMAT_ERROR"
	PlayerNotFoundError        = "PLAYER_NOT_FOUND_ERROR"
	PlayerAlreadyExistsError   = "PLAYER_ALREADY_EXISTS_ERROR"
	UsernameIsTakenError       = "USERNAME_IS_TAKEN_ERROR"

	UserNodeNotFoundError          = "USER_NODE_NOT_FOUND_ERROR"
	SelfFollowError                = "SELF_FOLLOW_ERROR"
	FollowershipAlreadyExistsError = "FOLLOWERSHIP_ALREADY_EXISTS_ERROR"
	FollowershipDoesNotExistError  = "FOLLOWERSHIP_DOES_NOT_EXIST_ERROR"

	InvalidCursorFormatError = "INVALID_CURSOR_FORMAT_ERROR"

	UnsupportedFollowSortType = "UNSUPPORTED_FOLLOW_SORT_TYPE"

	ShortSearchQueryError = "SHORT_SEARCH_QUERY_ERROR"

	RacquetProfileNotFoundError = "RACQUET_PROFILE_NOT_FOUND_ERROR"
)

func NewInvalidPlayerIDFormatError(id string) *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidPlayerIDFormatError,
		fmt.Sprintf("player id: %s must be a valid uuid", id),
	)
}

func NewPlayerNotFoundError(id uuid.UUID) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		PlayerNotFoundError,
		fmt.Sprintf("player with id: %s does not exist", id),
	)
}

func NewPlayerAlreadyExistsError(id uuid.UUID) *errors.APIError {
	return errors.NewAPIError(
		http.StatusConflict,
		PlayerAlreadyExistsError,
		fmt.Sprintf("player with id: %s already exists", id),
	)
}

func NewUsernameIsTakenError(username string) *errors.APIError {
	return errors.NewAPIError(
		http.StatusConflict,
		UsernameIsTakenError,
		fmt.Sprintf("username: %s has been taken already", username),
	)
}

func NewUserNodeNotFoundError(id string) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		UserNodeNotFoundError,
		fmt.Sprintf("user node with id: %s does not exist", id),
	)
}

func NewSelfFollowError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		SelfFollowError,
		fmt.Sprintf("you can't follow yourself"),
	)
}

func NewFollowershipAlreadyExistsError(id uuid.UUID) *errors.APIError {
	return errors.NewAPIError(
		http.StatusConflict,
		FollowershipAlreadyExistsError,
		fmt.Sprintf("you already follow user with id: %s", id),
	)
}

func NewFollowershipDoesNotExistError(id uuid.UUID) *errors.APIError {
	return errors.NewAPIError(
		http.StatusConflict,
		FollowershipDoesNotExistError,
		fmt.Sprintf("you do not follow user with id: %s", id),
	)
}

func NewInvalidCursorFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidCursorFormatError,
		"cursor must be an int64",
	)
}

func NewUnsupportedFollowSortType(sortType models.FollowSortType) *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		UnsupportedFollowSortType,
		fmt.Sprintf("unsupported follow sort type: %s", sortType),
	)
}

func NewShortSearchQueryError(minLength int) *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		ShortSearchQueryError,
		fmt.Sprintf("search query must contain at least %d characters", minLength),
	)
}

func NewRacquetProfileNotFoundError(id uuid.UUID, sportType domain.RacquetSportType) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		RacquetProfileNotFoundError,
		fmt.Sprintf("racquet profile with sport type: %d for player with id: %s does not exist", sportType, id),
	)
}
