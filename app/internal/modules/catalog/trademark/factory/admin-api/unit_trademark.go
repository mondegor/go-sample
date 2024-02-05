package factory

import (
	"context"
	module "go-sample/internal/modules/catalog/trademark"
	http_v1 "go-sample/internal/modules/catalog/trademark/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/trademark/entity/admin-api"
	"go-sample/internal/modules/catalog/trademark/factory"
	repository "go-sample/internal/modules/catalog/trademark/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/trademark/usecase/admin-api"

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

func newUnitTrademark(ctx context.Context, opts factory.Options) (*http_v1.Trademark, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Trademark{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewTrademarkPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewTrademark(storage, opts.EventEmitter, opts.UsecaseHelper)
	controller := http_v1.NewTrademark(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
