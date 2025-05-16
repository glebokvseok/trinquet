package authmw

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	auth "github.com/move-mates/trinquet/library/auth/pkg"
	"github.com/move-mates/trinquet/library/auth/pkg/errors"
	"github.com/move-mates/trinquet/library/auth/pkg/models"
	"go.uber.org/fx"
)

type UserAuthMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"user_auth_handler"`
}

func provideUserAuthMiddleware(config JWTConfig) UserAuthMiddleware {
	return UserAuthMiddleware{
		Middleware: echojwt.WithConfig(
			echojwt.Config{
				SigningKey: []byte(config.UserAccessTokenSigningKey),
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
