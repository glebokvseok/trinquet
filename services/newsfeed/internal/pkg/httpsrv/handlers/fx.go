package handlers

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"handlers",
		fx.Provide(
			provideCommentHandler,
			provideNewsfeedHandler,
			providePostEventHandler,
			providePostMediaHandler,
		),
	)
}

type Params struct {
	fx.In

	CommentHandler   *CommentHandler
	NewsfeedHandler  *NewsfeedHandler
	PostEventHandler *PostEventHandler
	PostMediaHandler *PostMediaHandler
}
