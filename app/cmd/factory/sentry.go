package factory

import (
	"context"
	"fmt"

	"github.com/mondegor/go-sample/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsentry"
)

// NewSentry - comment func.
func NewSentry(ctx context.Context, cfg config.Config) (*mrsentry.Adapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init sentry")

	client, err := mrsentry.New(
		mrsentry.Options{
			Dsn:              cfg.Sentry.Dsn,
			Environment:      cfg.AppEnvironment,
			TracesSampleRate: cfg.Sentry.TracesSampleRate,
			FlushTimeout:     cfg.Sentry.FlushTimeout,
			StackTraceBounds: cfg.Debugging.ErrorCaller.UpperBounds,
			IsDebug:          cfg.Debugging.Debug,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("sentry.Init: %w", err)
	}

	return client, nil
}
