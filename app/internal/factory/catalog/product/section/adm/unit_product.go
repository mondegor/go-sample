package adm

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/repository"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/product"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitProduct(ctx context.Context, opts product.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitProduct(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitProduct(ctx context.Context, opts product.Options) (*httpv1.Product, error) {
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
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}
