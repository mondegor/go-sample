package pub

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/repository"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/category"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCategory(ctx context.Context, opts category.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCategory(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCategory(_ context.Context, opts category.Options) (*httpv1.Category, error) { //nolint:unparam
	storage := repository.NewCategoryPostgres(
		opts.DBConnManager,
		mrpostgres.NewSQLBuilderSelect(
			mrpostgres.NewSQLBuilderWhere(),
			nil,
			mrpostgres.NewSQLBuilderLimit(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCategory(
		storage,
		opts.UsecaseHelper,
		opts.UnitCategory.ImageURLBuilder,
		opts.UnitCategory.Dictionary,
	)
	controller := httpv1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
