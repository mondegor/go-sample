package factory

import (
	"context"
	http_v1 "go-sample/internal/modules/catalog/category/controller/http_v1/public-api"
	"go-sample/internal/modules/catalog/category/factory"
	repository "go-sample/internal/modules/catalog/category/infrastructure/repository/public-api"
	usecase "go-sample/internal/modules/catalog/category/usecase/public-api"

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

func newUnitCategory(ctx context.Context, opts factory.Options) (*http_v1.Category, error) {
	storage := repository.NewCategoryPostgres(
		opts.PostgresAdapter,
		mrpostgres.NewSqlBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			nil,
			mrpostgres.NewSqlBuilderPager(opts.PageSizeMax),
		),
	)
	useCase := usecase.NewCategory(
		storage,
		opts.UsecaseHelper,
		opts.UnitCategory.ImageURLBuilder,
		opts.UnitCategory.Dictionary,
	)
	controller := http_v1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
	)

	return controller, nil
}
