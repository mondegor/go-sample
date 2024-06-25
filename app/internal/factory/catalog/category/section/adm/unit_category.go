package adm

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/repository"
	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/category"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCategory(ctx context.Context, opts category.Options) ([]mrserver.HttpController, error) {
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

func newUnitCategory(ctx context.Context, opts category.Options) (*httpv1.Category, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(ctx, entity.Category{})
	if err != nil {
		return nil, err
	}

	storage := repository.NewCategoryPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			mrpostgres.NewSQLBuilderOrderBy(ctx, metaOrderBy.DefaultSort()),
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCategory(storage, opts.EventEmitter, opts.UsecaseHelper, opts.UnitCategory.ImageURLBuilder)
	controller := httpv1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		metaOrderBy,
	)

	return controller, nil
}

func newUnitCategoryImage(_ context.Context, opts category.Options) (*httpv1.CategoryImage, error) { //nolint:unparam
	storage := repository.NewCategoryImagePostgres(opts.DBConnManager)
	useCase := usecase.NewCategoryImage(
		storage,
		opts.UnitCategory.ImageFileAPI,
		opts.Locker,
		opts.EventEmitter,
		opts.UsecaseHelper,
	)
	controller := httpv1.NewCategoryImage(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
