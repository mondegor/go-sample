package factory

import (
	"context"
	module "go-sample/internal/modules/catalog/category"
	http_v1 "go-sample/internal/modules/catalog/category/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/category/entity/admin-api"
	"go-sample/internal/modules/catalog/category/factory"
	repository "go-sample/internal/modules/catalog/category/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/category/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitCategory(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCategory(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	if c, err := newUnitCategoryImage(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCategory(ctx context.Context, opts factory.Options) (*http_v1.Category, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Category{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewCategoryPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCategory(storage, opts.EventEmitter, opts.UsecaseHelper, opts.UnitCategory.ImageURLBuilder)
	controller := http_v1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}

func newUnitCategoryImage(ctx context.Context, opts factory.Options) (*http_v1.CategoryImage, error) {
	storage := repository.NewCategoryImagePostgres(opts.PostgresAdapter)
	service := usecase.NewCategoryImage(
		storage,
		opts.UnitCategory.ImageFileAPI,
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := http_v1.NewCategoryImage(
		opts.RequestParser,
		mrresponse.NewFileSender(opts.ResponseSender),
		service,
	)

	return controller, nil
}
