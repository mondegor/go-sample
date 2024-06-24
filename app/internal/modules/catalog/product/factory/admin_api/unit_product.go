package factory

import (
	"context"

	httpv1 "go-sample/internal/modules/catalog/product/controller/httpv1/admin_api"
	entity "go-sample/internal/modules/catalog/product/entity/admin_api"
	"go-sample/internal/modules/catalog/product/factory"
	repository "go-sample/internal/modules/catalog/product/infrastructure/repository/admin_api"
	usecase "go-sample/internal/modules/catalog/product/usecase/admin_api"

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

func newUnitProduct(ctx context.Context, opts factory.Options) (*httpv1.Product, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Product{})
	if err != nil {
		return nil, err
	}

	entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(ctx, entity.Product{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewProductPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
		mrpostgres.NewSQLBuilderUpdateWithMeta(
			entityMetaUpdate,
			mrpostgres.NewSQLBuilderSet(),
			nil,
		),
	)
	useCase := usecase.NewProduct(
		storage,
		opts.CategoryAPI,
		opts.TrademarkAPI,
		opts.OrdererAPI,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := httpv1.NewProduct(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
