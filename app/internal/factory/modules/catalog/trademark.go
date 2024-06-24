package catalog

import (
	"context"

	"go-sample/pkg/modules/catalog/api"

	"go-sample/internal/app"
	view_shared "go-sample/internal/modules/catalog/trademark/controller/httpv1/shared/view"
	"go-sample/internal/modules/catalog/trademark/factory"
	factory_api "go-sample/internal/modules/catalog/trademark/factory/api"
	usecase_api "go-sample/internal/modules/catalog/trademark/usecase/api"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewTrademarkModuleOptions - comment func.
func NewTrademarkModuleOptions(_ context.Context, opts app.Options) (factory.Options, error) {
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
			// opts.RequestParsers.UUID,
			opts.RequestParsers.Validator,
			// opts.RequestParsers.File,
			// opts.RequestParsers.Image,
			opts.RequestParsers.ItemStatus,
		),
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewTrademarkAPI - comment func.
func NewTrademarkAPI(ctx context.Context, opts app.Options) (*usecase_api.Trademark, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog trademark API")

	return factory_api.NewTrademark(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterTrademarkErrors - comment func.
func RegisterTrademarkErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.TrademarkErrors()))
}
