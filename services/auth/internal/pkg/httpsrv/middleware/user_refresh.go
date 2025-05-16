package middleware

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	"github.com/move-mates/trinquet/library/auth/pkg/errors"
	"github.com/move-mates/trinquet/library/auth/pkg/models"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/generators"
	"go.uber.org/fx"
)

type UserRefreshMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"user_refresh_handler"`
}

func provideUserRefreshMiddleware(config generators.JWTConfig) UserRefreshMiddleware {
	return UserRefreshMiddleware{
		Middleware: echojwt.WithConfig(
			echojwt.Config{
				SigningKey: []byte(config.UserRefreshTokenSigningKey),
				ContextKey: auth.UserAuthDataContextKey,
				NewClaimsFunc: func(ctx echo.Context) gojwt.Claims {
					return new(models.CustomClaims)
				},
				Skipper: func(ctx echo.Context) bool {
					return false
				},
				ErrorHandler: func(ctx echo.Context, err error) error {
					if err != nil {
						return errors.NewInvalidOrExpiredJWTError()
					}

					return nil
				},
			},
		),
	}
}
