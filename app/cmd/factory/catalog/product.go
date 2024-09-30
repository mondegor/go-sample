package catalog

import (
	"context"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/product/module"

	"github.com/mondegor/go-sample/internal/factory/catalog/product"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewProductModuleOptions - создаёт объект product.Options.
func NewProductModuleOptions(_ context.Context, opts app.Options) (product.Options, error) {
	return product.Options{
		EventEmitter:  opts.EventEmitter,
		UseCaseHelper: opts.UseCaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: product.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		CategoryAPI:  opts.CatalogCategoryAvailabilityAPI,
		OrdererAPI:   opts.OrdererAPI,
		TrademarkAPI: opts.CatalogTrademarkAvailabilityAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// RegisterProductErrors - comment func.
func RegisterProductErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
