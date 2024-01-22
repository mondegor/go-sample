package factory

import (
	module "go-sample/internal/modules/catalog"
	http_v1 "go-sample/internal/modules/catalog/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	"go-sample/internal/modules/catalog/factory"
	repository "go-sample/internal/modules/catalog/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func createUnitCategory(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitCategory(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	if c, err := newUnitCategoryImage(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitCategory(opts *factory.Options) (*http_v1.Category, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Category{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewCategoryPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCategory(storage, opts.EventBox, opts.ServiceHelper, opts.UnitCategory.ImageURLBuilder)
	controller := http_v1.NewCategory(
		opts.RequestParser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}

func newUnitCategoryImage(opts *factory.Options) (*http_v1.CategoryImage, error) {
	storage := repository.NewCategoryImagePostgres(opts.PostgresAdapter)
	service := usecase.NewCategoryImage(
		storage,
		opts.UnitCategory.ImageFileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)
	controller := http_v1.NewCategoryImage(
		opts.RequestParser,
		mrresponse.NewFileSender(opts.ResponseSender),
		service,
	)

	return controller, nil
}
