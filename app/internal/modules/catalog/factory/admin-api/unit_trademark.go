package factory

import (
	"go-sample/internal/modules"
	http_v1 "go-sample/internal/modules/catalog/controller/http_v1/admin-api"
	entity "go-sample/internal/modules/catalog/entity/admin-api"
	repository "go-sample/internal/modules/catalog/infrastructure/repository/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	unitNameTrademark = moduleName + ".Trademark"
)

func newUnitTrademark(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
	storage *repository.Trademark,
	listSorter mrview.ListSorter,
) error {
	service := usecase.NewTrademark(storage, opts.EventBox, opts.ServiceHelper)
	*c = append(*c, http_v1.NewTrademark(section, service, listSorter))

	return nil
}

func newUnitTrademarkStorage(opts *modules.Options) (*repository.Trademark, *mrsql.EntityMetaOrderBy, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Trademark{})

	if err != nil {
		return nil, nil, err
	}

	return repository.NewTrademark(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(1000),
		),
	), metaOrderBy, nil
}
