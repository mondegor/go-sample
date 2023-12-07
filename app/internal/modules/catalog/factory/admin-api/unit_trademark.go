package factory

import (
	"go-sample/internal/modules"
	module "go-sample/internal/modules/catalog"
	http_v1 "go-sample/internal/modules/catalog/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	repository "go-sample/internal/modules/catalog/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	usecase_shared "go-sample/internal/modules/catalog/usecase/shared"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

const (
	unitNameTrademark = moduleName + ".Trademark"
)

func newUnitTrademark(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
) (*usecase_shared.TrademarkAPI, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Trademark{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewTrademark(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewTrademark(storage, opts.EventBox, opts.ServiceHelper)
	serviceAPI := usecase_shared.NewTrademarkAPI(storage, opts.ServiceHelper)
	*c = append(*c, http_v1.NewTrademark(section, service, metaOrderBy))

	return serviceAPI, nil
}
