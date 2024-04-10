package factory

import (
	"context"
	http_v1 "go-sample/internal/modules/catalog/category/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/category/entity/admin-api"
	"go-sample/internal/modules/catalog/category/factory"
	repository "go-sample/internal/modules/catalog/category/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/category/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
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
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCategory(storage, opts.EventEmitter, opts.UsecaseHelper, opts.UnitCategory.ImageURLBuilder)
	controller := http_v1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}

func newUnitCategoryImage(ctx context.Context, opts factory.Options) (*http_v1.CategoryImage, error) {
	storage := repository.NewCategoryImagePostgres(opts.PostgresAdapter)
	useCase := usecase.NewCategoryImage(
		storage,
		opts.UnitCategory.ImageFileAPI,
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := http_v1.NewCategoryImage(
		opts.RequestParser,
		mrresp.NewFileSender(opts.ResponseSender),
		useCase,
	)

	return controller, nil
}
