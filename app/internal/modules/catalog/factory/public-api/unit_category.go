package factory

import (
	module "go-sample/internal/modules/catalog"
	http_v1 "go-sample/internal/modules/catalog/controller/http_v1/public-api"
	"go-sample/internal/modules/catalog/factory"
	repository "go-sample/internal/modules/catalog/infrastructure/repository/public-api"
	usecase "go-sample/internal/modules/catalog/usecase/public-api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrcore"
)

func newUnitCategory(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	storage := repository.NewCategoryPostgres(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			nil,
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCategoryLangDecorator(
		usecase.NewCategory(storage, opts.ServiceHelper, opts.UnitCategory.ImageURLBuilder),
		opts.UnitCategory.Dictionary,
	)
	*c = append(*c, http_v1.NewCategory(section, service))

	return nil
}
