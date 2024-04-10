package factory

import (
	"context"
	"go-sample/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

func NewResponseSender(ctx context.Context, cfg config.Config) (*mrresp.Sender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base response sender")

	return mrresp.NewSender(mrjson.NewEncoder()), nil
}

func NewErrorResponseSender(ctx context.Context, cfg config.Config) (*mrresp.ErrorSender, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init error response sender")

	return mrresp.NewErrorSender(mrjson.NewEncoder()), nil
}
