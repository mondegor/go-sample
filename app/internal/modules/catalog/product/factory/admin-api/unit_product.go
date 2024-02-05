package factory

import (
	"context"
	module "go-sample/internal/modules/catalog/product"
	http_v1 "go-sample/internal/modules/catalog/product/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/product/entity/admin-api"
	"go-sample/internal/modules/catalog/product/factory"
	repository "go-sample/internal/modules/catalog/product/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/product/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitProduct(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitProduct(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitProduct(ctx context.Context, opts factory.Options) (*http_v1.Product, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Product{})

	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Product{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewProductPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
		mrsql.NewBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSqlBuilderSet(),
		),
	)
	service := usecase.NewProduct(
		storage,
		opts.CategoryAPI,
		opts.TrademarkAPI,
		opts.OrdererAPI,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := http_v1.NewProduct(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
