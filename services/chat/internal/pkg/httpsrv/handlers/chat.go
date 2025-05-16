package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	cmnerrors "github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/httpsrv/requests"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/httpsrv/responses"
	"github.com/move-mates/trinquet/services/chat/internal/pkg/managers"
	"net/http"
)

type ChatHandler struct {
	chatManager managers.ChatManager
	validator   *validator.Validate
}

func provideChatHandler(
	chatManager managers.ChatManager,
	validator *validator.Validate,
) *ChatHandler {
	return &ChatHandler{
		chatManager: chatManager,
		validator:   validator,
	}
}

func (handler *ChatHandler) GetChatMessages(ctx echo.Context) error {
	var request requests.GetChatMessagesRequest
	err := ctx.Bind(&request)
	if err != nil {
		return cmnerrors.NewRequestBodyParsingError(err)
	}

	err = handler.validator.Struct(request)
	if err != nil {
		return cmnerrors.NewInvalidRequestBodyFormatError(err)
	}

	messages, err := handler.chatManager.GetChatMessages(
		ctx.Request().Context(),
		auth.GetUserID(ctx),
		request.Chat,
		request.Since,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		responses.GetChatMessagesResponse{
			Messages: messages,
		},
	)
}
