package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/factory/catalog/product"
)

// NewProductModuleOptions - создаёт объект product.Options.
func NewProductModuleOptions(_ context.Context, opts app.Options) (product.Options, error) {
	return product.Options{
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		RequestParsers: product.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		CategoryAPI:  opts.CatalogCategoryAvailabilityAPI,
		TrademarkAPI: opts.CatalogTrademarkAvailabilityAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}
