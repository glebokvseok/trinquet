package psql

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func provideDatabase(config databaseConfig, lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.ConnectionString), &gorm.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.MaxConnectionLifetime)
	sqlDB.SetConnMaxIdleTime(config.MaxConnectionIdleTime)

	// TODO: разобраться с корректным освобождением ресурсов при завершении, хук сейчас не срабатывает
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
	})

	return db, nil
}
