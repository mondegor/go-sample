package factory

import (
    "go-sample/internal/factory"
    "go-sample/internal/modules"
    http_v1 "go-sample/internal/modules/catalog/controller/http_v1/public-api"
    repository "go-sample/internal/modules/catalog/infrastructure/repository/public-api"
    usecase "go-sample/internal/modules/catalog/usecase/public-api"

    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-storage/mrsql"
    "github.com/mondegor/go-webcore/mrcore"
)

const (
    unitNameCategory = "Category"
)

func newUnitCategory(
    c *[]mrcore.HttpController,
    opts *modules.Options,
    section mrcore.ClientSection,
) error {
    storage := repository.NewCategory(
        opts.PostgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgres.NewSqlBuilderWhere(),
            nil,
            mrpostgres.NewSqlBuilderPager(1000),
        ),
    )
    imagesURL := factory.NewBuilderImagesURL(opts.Cfg)
    service := usecase.NewCategory(storage, opts.ServiceHelper)
    *c = append(*c, http_v1.NewCategory(section, service, imagesURL))

    return nil
}
