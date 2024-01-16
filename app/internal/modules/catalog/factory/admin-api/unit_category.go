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

func newUnitCategory(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	if err := newUnitCategoryImage(c, opts, section); err != nil {
		return err
	}

	return newUnitCategoryMain(c, opts, section)
}

func newUnitCategoryMain(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Category{})

	if err != nil {
		return err
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
	*c = append(*c, http_v1.NewCategory(section, service, metaOrderBy))

	return nil
}

func newUnitCategoryImage(
	c *[]mrcore.HttpController,
	opts *factory.Options,
	section mrcore.ClientSection,
) error {
	storage := repository.NewCategoryImagePostgres(opts.PostgresAdapter)
	service := usecase.NewCategoryImage(
		storage,
		opts.UnitCategory.ImageFileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)
	*c = append(*c, http_v1.NewCategoryImage(section, service))

	return nil
}
