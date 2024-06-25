package catalog

import (
	"context"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/catalog/product/shared/validate"
	"github.com/mondegor/go-sample/internal/factory/catalog/product"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewProductModuleOptions - comment func.
func NewProductModuleOptions(_ context.Context, opts app.Options) (product.Options, error) {
	return product.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParser: validate.NewParser(
			// opts.RequestParsers.Bool,
			// opts.RequestParsers.DateTime,
			opts.RequestParsers.Int64,
			opts.RequestParsers.KeyInt32,
			opts.RequestParsers.ListSorter,
			opts.RequestParsers.ListPager,
			opts.RequestParsers.String,
			opts.RequestParsers.UUID,
			opts.RequestParsers.Validator,
			// opts.RequestParsers.File,
			// opts.RequestParsers.Image,
			opts.RequestParsers.ItemStatus,
		),
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
