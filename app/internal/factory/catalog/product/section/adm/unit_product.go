package adm

import (
	"context"

	"github.com/mondegor/go-components/factory/mrordering"
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/repository"
	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/product"
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
	entityMeta, err := mrsql.ParseEntity(mrlog.Ctx(ctx), entity.Product{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewProductPostgres(
		opts.DBConnManager,
		builder.NewSQL(
			builder.WithSQLSetMetaEntity(entityMeta.MetaUpdate()),
			builder.WithSQLOrderByDefaultSort(entityMeta.MetaOrderBy().DefaultSort()),
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewProduct(
		storage,
		opts.CategoryAPI,
		opts.TrademarkAPI,
		mrordering.NewComponentMover(
			opts.DBConnManager,
			mrsql.DBTableInfo{
				Name:       module.DBTableNameProducts,
				PrimaryKey: "product_id",
			},
			opts.EventEmitter,
		),
		opts.EventEmitter,
		opts.UseCaseErrorWrapper,
	)
	controller := httpv1.NewProduct(
		opts.RequestParsers.ExtendParser,
		opts.ResponseSender,
		useCase,
		entityMeta.MetaOrderBy(),
	)

	return controller, nil
}
