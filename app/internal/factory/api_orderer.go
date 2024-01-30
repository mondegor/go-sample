package factory

import (
	"context"
	"go-sample/internal/modules"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-webcore/mrlog"
)

func NewOrdererAPI(ctx context.Context, opts modules.Options) *mrorderer.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init orderer component")
	itemOrdererStorage := mrorderer.NewRepository(opts.PostgresAdapter)

	return mrorderer.NewComponent(itemOrdererStorage, opts.EventEmitter)
}
