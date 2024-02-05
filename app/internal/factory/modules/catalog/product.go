package factory_catalog

import (
	"context"
	"go-sample/internal"
	view_shared "go-sample/internal/modules/catalog/product/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/product/factory"
)

func NewProductModuleOptions(ctx context.Context, opts app.Options) (factory.Options, error) {
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

		CategoryAPI:  opts.CatalogCategoryAPI,
		OrdererAPI:   opts.OrdererAPI,
		TrademarkAPI: opts.CatalogTrademarkAPI,
	}, nil
}
