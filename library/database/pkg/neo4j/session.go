package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pkg/errors"
)

const (
	errorMessage = "unhandled panic occurred while executing neo4j request: %v"
)

type Result = neo4j.ResultWithContext

func NewReadSession(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	config SessionConfig,
) (neo4j.SessionWithContext, context.Context, context.CancelFunc) {
	return newSession(ctx, driver, config, neo4j.AccessModeRead)
}

func NewWriteSession(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	config SessionConfig,
) (neo4j.SessionWithContext, context.Context, context.CancelFunc) {
	return newSession(ctx, driver, config, neo4j.AccessModeWrite)
}

func newSession(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	config SessionConfig,
	accessMode neo4j.AccessMode,
) (neo4j.SessionWithContext, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, config.RequestTimeout)
	session := driver.NewSession(
		ctx,
		neo4j.SessionConfig{
			DatabaseName: config.Database,
			AccessMode:   accessMode,
		},
	)

	return session, ctx, cancel
}

func SessionFinalizer(
	ctx context.Context,
	session neo4j.SessionWithContext,
	returnErr *error,
) {
	if r := recover(); r != nil {
		if returnErr != nil {
			*returnErr = errors.Errorf(errorMessage, r)
		}
	}

	if session != nil {
		// TODO: надо этот кейс как-нибудь обрабатывать, например, логировать
		_ = session.Close(ctx)
	}
}
