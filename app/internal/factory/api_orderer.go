package factory

import (
	"context"

	"go-sample/internal/app"

	sortfactory "github.com/mondegor/go-components/mrsort/factory"
	"github.com/mondegor/go-components/mrsort/orderer"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewOrdererAPI - comment func.
func NewOrdererAPI(ctx context.Context, opts app.Options) *orderer.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init orderer component")

	return sortfactory.NewComponentOrderer(
		sortfactory.ComponentOptions{
			DBClient:     opts.PostgresConnManager,
			EventEmitter: opts.EventEmitter,
		},
	)
}
