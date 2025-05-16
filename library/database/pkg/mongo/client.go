package mongo

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"time"
)

type Client = *mongo.Client

func provideClient(config databaseConfig, lc fx.Lifecycle) (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opts := options.Client().
		ApplyURI(config.ConnectionString).
		SetMaxPoolSize(config.MaxConnectionPoolSize).
		SetMaxConnIdleTime(config.MaxConnectionIdleTime)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: разобраться с корректным освобождением ресурсов при завершении, хук сейчас не срабатывает
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return client, nil
}
