package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/repository"
	"github.com/mondegor/go-sample/internal/catalog/trademark/section/adm/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark"
)

func createUnitTrademark(ctx context.Context, opts trademark.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitTrademark(ctx context.Context, opts trademark.Options) (*httpv1.Trademark, error) {
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.Trademark{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewTrademarkPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewTrademark(storage, opts.EventEmitter, opts.UseCaseErrorWrapper)
	controller := httpv1.NewTrademark(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
