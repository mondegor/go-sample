package factory

import (
	"context"
	module "go-sample/internal/modules/catalog"
	http_v1 "go-sample/internal/modules/catalog/controller/http_v1/public-api"
	"go-sample/internal/modules/catalog/factory"
	repository "go-sample/internal/modules/catalog/infrastructure/repository/public-api"
	usecase "go-sample/internal/modules/catalog/usecase/public-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
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
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			nil,
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCategoryLangDecorator(
		usecase.NewCategory(storage, opts.UsecaseHelper, opts.UnitCategory.ImageURLBuilder),
		opts.UnitCategory.Dictionary,
	)
	controller := http_v1.NewCategory(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		service,
	)

	return controller, nil
}
