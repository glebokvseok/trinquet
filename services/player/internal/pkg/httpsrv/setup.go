package httpsrv

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	handlerspkg "github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/handlers"
	middlewarepkg "github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/middleware"
)

const (
	baseGroupPrefix = "/v0/player-service"

	avatarGroupPrefix = "/avatar"
	playerGroupPrefix = "/player"

	racquetGroupPrefix = "/racquet"
	profileGroupPrefix = "/profile"
	matchGroupPrefix   = "/match"
)

func ProvideHttpServerSetupFunc(
	handlers handlerspkg.Params,
	middleware middlewarepkg.Params,
) httpsrv.HttpServerSetupFunc {
	return func(e *echo.Echo) {
		registerRoutes(e, handlers, middleware)
	}
}

func registerRoutes(
	e *echo.Echo,
	handlers handlerspkg.Params,
	middleware middlewarepkg.Params,
) {
	base := e.Group(baseGroupPrefix)
	base.Use(middleware.PanicHandler)
	base.Use(middleware.ErrorHandler)
	base.Use(middleware.CORS)
	base.GET(httpsrv.HealthCheckPath, httpsrv.HealthCheckFunc)
	base.Use(middleware.SignatureHandler)
	base.Use(middleware.RequestIDHandler)
	base.Use(middleware.UserAuthHandler)

	player := base.Group(playerGroupPrefix)

	player.GET("/self", handlers.PlayerHandler.GetSelf)
	player.GET(fmt.Sprintf("/:%s", handlerspkg.PlayerIDParamName), handlers.PlayerHandler.GetOther)
	player.GET("/search", handlers.PlayerHandler.SearchPlayer)
	player.POST("/create", handlers.PlayerHandler.CreatePlayer)
	player.PUT("/update", handlers.PlayerHandler.UpdatePlayer)

	player.POST("/follow", handlers.PlayerRelationHandler.FollowPlayer)
	player.POST("/unfollow", handlers.PlayerRelationHandler.UnfollowPlayer)
	player.GET("/followers", handlers.PlayerRelationHandler.GetFollowers)
	player.GET("/following", handlers.PlayerRelationHandler.GetFollowing)

	player.GET("/achievements", handlers.PlayerHandler.GetAchievements)

	avatar := base.Group(avatarGroupPrefix)

	avatar.GET("/download/self", handlers.AvatarHandler.DownloadSelf)
	avatar.GET(fmt.Sprintf("/download/:%s", handlerspkg.PlayerIDParamName), handlers.AvatarHandler.DownloadOther)
	avatar.POST("/upload", handlers.AvatarHandler.Upload)

	racquet := base.Group(racquetGroupPrefix)

	profile := racquet.Group(profileGroupPrefix)
	profile.PATCH("/update", handlers.RacquetProfileHandler.UpdateRacquetProfiles)

	match := racquet.Group(matchGroupPrefix)
	match.POST("/create", handlers.RacquetMatchHandler.CreateMatch)
	match.POST("/update", handlers.RacquetMatchHandler.UpdateMatch)
	match.POST("/start", handlers.RacquetMatchHandler.StartMatch)
	match.POST("/submit-results", handlers.RacquetMatchHandler.SubmitResults)
	match.POST("/invite", handlers.RacquetMatchHandler.InvitePlayer)
	match.POST("/join", handlers.RacquetMatchHandler.JoinMatch)
	match.GET("/history", handlers.RacquetMatchHandler.GetMatchHistory)
}
