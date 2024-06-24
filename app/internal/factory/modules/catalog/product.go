package catalog

import (
	"context"

	"go-sample/internal/app"
	view_shared "go-sample/internal/modules/catalog/product/controller/httpv1/shared/view"
	"go-sample/internal/modules/catalog/product/factory"
	"go-sample/internal/modules/catalog/product/module"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
)

// NewProductModuleOptions - comment func.
func NewProductModuleOptions(_ context.Context, opts app.Options) (factory.Options, error) {
	return factory.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParser: view_shared.NewParser(
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

		CategoryAPI:  opts.CatalogCategoryAPI,
		OrdererAPI:   opts.OrdererAPI,
		TrademarkAPI: opts.CatalogTrademarkAPI,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// RegisterProductErrors - comment func.
func RegisterProductErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
