package factory

import (
	"context"

	"go-sample/internal/app"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-storage/mrstorage/mrprometheus"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewPostgres - comment func.
func NewPostgres(ctx context.Context, opts app.Options) (*mrpostgres.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init postgres connection")

	cfg := opts.Cfg
	postgresOpts := mrpostgres.Options{
		Host:         cfg.Storage.Host,
		Port:         cfg.Storage.Port,
		Username:     cfg.Storage.Username,
		Password:     cfg.Storage.Password,
		Database:     cfg.Storage.Database,
		MaxPoolSize:  cfg.Storage.MaxPoolSize,
		ConnAttempts: 1,
		ConnTimeout:  cfg.Storage.Timeout,
	}

	conn := mrpostgres.New()

	if err := conn.Connect(ctx, postgresOpts); err != nil {
		return nil, err
	}

	opts.Prometheus.MustRegister(
		mrprometheus.NewDBCollector(
			"pgx",
			func() mrstorage.DBStatProvider {
				return conn.Cli().Stat()
			},
			map[string]string{
				"db_name": cfg.Storage.Database,
			},
		),
	)

	return conn, conn.Ping(ctx)
}

// NewPostgresConnManager - comment func.
func NewPostgresConnManager(_ context.Context, conn *mrpostgres.ConnAdapter) *mrpostgres.ConnManager {
	return mrpostgres.NewConnManager(conn)
}
