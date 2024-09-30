package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark/api/availability"
	"github.com/mondegor/go-sample/pkg/catalog/api"
)

// NewTrademarkModuleOptions - создаёт объект trademark.Options.
func NewTrademarkModuleOptions(_ context.Context, opts app.Options) (trademark.Options, error) {
	return trademark.Options{
		EventEmitter:  opts.EventEmitter,
		UseCaseHelper: opts.UseCaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		RequestParsers: trademark.RequestParsers{
			Parser:       opts.RequestParsers.Parser,
			ExtendParser: opts.RequestParsers.ExtendParser,
		},
		ResponseSender: opts.ResponseSenders.Sender,

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewTrademarkAvailabilityAPI - создаёт объект usecase.Trademark.
func NewTrademarkAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.Trademark, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog trademark availability API")

	return availability.NewTrademark(opts.PostgresConnManager, opts.UseCaseErrorWrapper), nil
}

// RegisterTrademarkErrors - comment func.
func RegisterTrademarkErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.TrademarkErrors()))
}
