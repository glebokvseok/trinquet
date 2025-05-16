package neo4j

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"io"
	"os"
)

type Driver = neo4j.DriverWithContext

func provideDriver(conf databaseConfig, lc fx.Lifecycle) (driver Driver, returnErr error) {
	var tlsConfig *tls.Config
	if extensions.IsNotEmpty(conf.SSLRootCert) {
		certPool := x509.NewCertPool()

		cert, err := readCert(conf.SSLRootCert)
		if err != nil {
			return nil, err
		}

		ok := certPool.AppendCertsFromPEM(cert)
		if !ok {
			return nil, errors.Errorf("failed to append root cert to pool")
		}

		tlsConfig = &tls.Config{
			RootCAs: certPool,
		}
	}

	driver, err := neo4j.NewDriverWithContext(
		conf.URI,
		neo4j.BasicAuth(conf.Username, conf.Password, ""),
		func(driverConfig *config.Config) {
			driverConfig.TlsConfig = tlsConfig
			driverConfig.MaxConnectionPoolSize = conf.MaxConnectionPoolSize
			driverConfig.MaxConnectionLifetime = conf.MaxConnectionLifetime
		},
	)

	// TODO: разобраться с корректным освобождением ресурсов при завершении, хук сейчас не срабатывает
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return driver.Close(ctx)
		},
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return driver, nil
}

func readCert(sslRootCert string) (data []byte, returnErr error) {
	file, err := os.Open(sslRootCert)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			returnErr = errors.WithStack(err)
		}
	}()

	data, err = io.ReadAll(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}
