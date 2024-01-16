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
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitTrademark(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Trademark{})

	if err != nil {
		return err
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
	*c = append(*c, http_v1.NewTrademark(section, service, metaOrderBy))

	return nil
}
