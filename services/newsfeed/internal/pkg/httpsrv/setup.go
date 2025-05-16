package httpsrv

import (
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/httpsrv"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/httpsrv/handlers"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/httpsrv/middleware"
)

const (
	baseGroupPrefix = "/v0/newsfeed-service"

	newsfeedGroupPrefix = "/newsfeed"
	postGroupPrefix     = "/post"
)

func ProvideHttpServerSetupFunc(
	handlers handlers.Params,
	middleware middleware.Params,
) httpsrv.HttpServerSetupFunc {
	return func(e *echo.Echo) {
		registerRoutes(e, handlers, middleware)
	}
}

func registerRoutes(
	e *echo.Echo,
	handlers handlers.Params,
	middleware middleware.Params,
) {
	base := e.Group(baseGroupPrefix)
	base.Use(middleware.PanicHandler)
	base.Use(middleware.ErrorHandler)
	base.Use(middleware.CORS)
	base.GET(httpsrv.HealthCheckPath, httpsrv.HealthCheckFunc)
	base.Use(middleware.SignatureHandler)
	base.Use(middleware.RequestIDHandler)

	post := base.Group(postGroupPrefix)
	post.Use(middleware.UserAuthHandler)
	post.GET("/comments", handlers.CommentHandler.GetCommentSection)
	post.POST("/create", handlers.PostEventHandler.CreatePost)
	post.POST("/like", handlers.PostEventHandler.LikePost)
	post.POST("/unlike", handlers.PostEventHandler.UnlikePost)
	post.POST("/comment", handlers.PostEventHandler.CommentPost)
	post.POST("/comment/reply", handlers.PostEventHandler.ReplyPostComment)
	post.POST("/media/upload", handlers.PostMediaHandler.UploadPostMedia)

	newsfeed := base.Group(newsfeedGroupPrefix)
	newsfeed.Use(middleware.UserAuthHandler)
	newsfeed.GET("", handlers.NewsfeedHandler.GetNewsfeed)
}
