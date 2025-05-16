package job

import "context"

type BackgroundJob interface {
	Run(context.Context) error
}
