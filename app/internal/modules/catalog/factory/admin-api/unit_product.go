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
)

const (
    unitNameProduct = moduleName + ".Product"
)

func newUnitProduct(
    c *[]mrcore.HttpController,
    opts *modules.Options,
    section mrcore.ClientSection,
    categoryStorage usecase.CategoryStorage,
    trademarkStorage usecase.TrademarkStorage,
) error {
    metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Product{})

    if err != nil {
        return err
    }

    entityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity.Product{})

    if err != nil {
        return err
    }

    storage := repository.NewProduct(
        opts.PostgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgres.NewSqlBuilderWhere(),
            mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
            mrpostgres.NewSqlBuilderPager(1000),
        ),
        mrsql.NewBuilderUpdateWithMeta(
            entityMetaUpdate,
            mrpostgres.NewSqlBuilderSet(),
        ),
    )
    service := usecase.NewProduct(opts.OrdererComponent, storage, categoryStorage, trademarkStorage, opts.EventBox, opts.ServiceHelper)
    *c = append(*c, http_v1.NewProduct(section, service, metaOrderBy))

    return nil
}
