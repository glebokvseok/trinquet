package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/managers"
	"net/http"
)

type RegistrationHandler struct {
	registrationManager managers.RegistrationManager
	validator           *validator.Validate
}

func provideRegistrationHandler(
	registrationManager managers.RegistrationManager,
	validator *validator.Validate,
) *RegistrationHandler {
	return &RegistrationHandler{
		registrationManager: registrationManager,
		validator:           validator,
	}
}

func (handler *RegistrationHandler) RegisterUser(ctx echo.Context) error {
	var request requests.RegisterUserRequest
	err := ctx.Bind(&request)
	if err != nil {
		return errors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return errors.NewInvalidRequestBodyFormatError(err)
	}

	accessToken, refreshToken, err := handler.registrationManager.RegisterUser(
		ctx.Request().Context(),
		request.Email,
		request.Password,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.RegisterUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
}
