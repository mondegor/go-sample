package catalog

import (
	"context"

	"github.com/mondegor/go-sample/pkg/catalog/api"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/usecase"
	"github.com/mondegor/go-sample/internal/catalog/trademark/shared/validate"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark/api/availability"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewTrademarkModuleOptions - comment func.
func NewTrademarkModuleOptions(_ context.Context, opts app.Options) (trademark.Options, error) {
	return trademark.Options{
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
func NewTrademarkAPI(ctx context.Context, opts app.Options) (*usecase.Trademark, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog trademark API")

	return availability.NewTrademark(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterTrademarkErrors - comment func.
func RegisterTrademarkErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.TrademarkErrors()))
}
