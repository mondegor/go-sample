package factory

import (
	"context"

	httpv1 "go-sample/internal/modules/catalog/category/controller/httpv1/public_api"
	"go-sample/internal/modules/catalog/category/factory"
	repository "go-sample/internal/modules/catalog/category/infrastructure/repository/public_api"
	usecase "go-sample/internal/modules/catalog/category/usecase/public_api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitCategory(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCategory(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCategory(_ context.Context, opts factory.Options) (*httpv1.Category, error) { //nolint:unparam
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
