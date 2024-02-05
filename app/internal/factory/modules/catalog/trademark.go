package factory_catalog

import (
	"context"
	"go-sample/internal"
	view_shared "go-sample/internal/modules/catalog/trademark/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/trademark/factory"
	factory_api "go-sample/internal/modules/catalog/trademark/factory/api"
	usecase_api "go-sample/internal/modules/catalog/trademark/usecase/api"

	"github.com/mondegor/go-webcore/mrlog"
)

func NewTrademarkModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
	return factory.Options{
		EventEmitter:    opts.EventEmitter,
		UsecaseHelper:   opts.UsecaseHelper,
		PostgresAdapter: opts.PostgresAdapter,
		RequestParser: view_shared.NewParser(
			// opts.RequestParsers.Bool,
			// opts.RequestParsers.DateTime,
			opts.RequestParsers.Int64,
			opts.RequestParsers.KeyInt32,
			opts.RequestParsers.SortPage,
			opts.RequestParsers.String,
			// opts.RequestParsers.UUID,
			opts.RequestParsers.Validator,
			// opts.RequestParsers.File,
			// opts.RequestParsers.Image,
			opts.RequestParsers.ItemStatus,
		),
		ResponseSender: opts.ResponseSender,
	}, nil
}

func NewTrademarkAPI(ctx context.Context, opts app.Options) (*usecase_api.Trademark, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog trademark API")

	return factory_api.NewTrademark(opts.PostgresAdapter, opts.UsecaseHelper), nil
}
