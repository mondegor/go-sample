package factory

import (
	"context"

	"go-sample/config"
	"go-sample/internal/app"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// CreateResponseSenders - comment func.
func CreateResponseSenders(ctx context.Context, _ config.Config) (app.ResponseSenders, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base response senders")

	sender := mrresp.NewSender(mrjson.NewEncoder())

	return app.ResponseSenders{
		Sender:     mrresp.NewSender(mrjson.NewEncoder()),
		FileSender: mrresp.NewFileSender(sender),
	}, nil
}

// NewErrorResponseSender - comment func.
func NewErrorResponseSender(ctx context.Context, opts app.Options) (*mrresp.ErrorSender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init error response sender")

	return mrresp.NewErrorSender(
		mrjson.NewEncoder(),
		opts.ErrorHandler,
		mrserver.NewHttpErrorStatusGetter(
			opts.Cfg.Debugging.UnexpectedHttpStatus,
		),
		opts.Cfg.Debugging.UnexpectedHttpStatus,
		opts.Cfg.Debugging.Debug,
	), nil
}
