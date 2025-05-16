package async

import (
	"context"
	"github.com/move-mates/trinquet/library/common/pkg/result"
	"github.com/pkg/errors"
)

func SafeGoRes[TRes any](
	ctx context.Context,
	ch chan<- result.Result[TRes],
	fn func(ctx context.Context) (TRes, error),
) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				select {
				case ch <- result.Result[TRes]{
					Err: errors.Errorf("unhandeled panic occured: %+v", r),
				}:
				case <-ctx.Done():
				}
			}

			return
		}()

		data, err := fn(ctx)
		select {
		case ch <- result.Result[TRes]{Data: data, Err: err}:
		case <-ctx.Done():
		}
	}()
}
