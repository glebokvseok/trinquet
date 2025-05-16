package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/managers"
	"net/http"
	"strings"
)

type AuthenticationHandler struct {
	authenticationManager managers.AuthenticationManager
	validator             *validator.Validate
}

func provideAuthenticationHandler(
	authenticationManager managers.AuthenticationManager,
	validator *validator.Validate,
) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticationManager: authenticationManager,
		validator:             validator,
	}
}

func (handler *AuthenticationHandler) AuthenticateUser(ctx echo.Context) error {
	var request requests.AuthenticateUserRequest
	err := ctx.Bind(&request)
	if err != nil {
		return errors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return errors.NewInvalidRequestBodyFormatError(err)
	}

	accessToken, refreshToken, err := handler.authenticationManager.AuthenticateUser(
		ctx.Request().Context(),
		request.Email,
		request.Password,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.AuthenticateUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
}

func (handler *AuthenticationHandler) RefreshUser(ctx echo.Context) error {
	refreshToken := strings.TrimPrefix(ctx.Request().Header.Get("Authorization"), "Bearer ")

	accessToken, newRefreshToken, err := handler.authenticationManager.RefreshUser(refreshToken)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.RefreshUserResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		},
	)
}
