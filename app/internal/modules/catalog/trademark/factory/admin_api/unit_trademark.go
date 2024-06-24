package factory

import (
	"context"

	httpv1 "go-sample/internal/modules/catalog/trademark/controller/httpv1/admin_api"
	entity "go-sample/internal/modules/catalog/trademark/entity/admin_api"
	"go-sample/internal/modules/catalog/trademark/factory"
	repository "go-sample/internal/modules/catalog/trademark/infrastructure/repository/admin_api"
	usecase "go-sample/internal/modules/catalog/trademark/usecase/admin_api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitTrademark(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitTrademark(ctx context.Context, opts factory.Options) (*httpv1.Trademark, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Trademark{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewTrademarkPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewTrademark(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := httpv1.NewTrademark(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
