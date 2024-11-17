package pub

import (
	"context"

	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/controller/httpv1"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/repository"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/usecase"
	"github.com/mondegor/go-sample/internal/factory/catalog/category"
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
		builder.NewSQL(
			builder.WithSQLLimitMaxSize(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCategory(
		storage,
		opts.UseCaseErrorWrapper,
		opts.UnitCategory.ImageURLBuilder,
		opts.UnitCategory.Dictionary,
	)
	controller := httpv1.NewCategory(
		opts.RequestParsers.ModuleParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
