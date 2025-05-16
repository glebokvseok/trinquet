package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/auth/pkg/models"
)

func GetUserID(ctx echo.Context) uuid.UUID {
	return ctx.Get(UserAuthDataContextKey).(*jwt.Token).Claims.(*models.CustomClaims).UserId
}
