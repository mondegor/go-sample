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
)

func createUnitTrademark(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitTrademark(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitTrademark(opts *factory.Options) (*http_v1.Trademark, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Trademark{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewTrademarkPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewTrademark(storage, opts.EventBox, opts.ServiceHelper)
	controller := http_v1.NewTrademark(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		service,
		metaOrderBy,
	)

	return controller, nil
}
