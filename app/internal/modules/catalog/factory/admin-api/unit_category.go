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
	unitNameCategory = moduleName + ".Category"
)

func newUnitCategory(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
	serviceImage *usecase.CategoryImage,
) (*usecase_shared.CategoryAPI, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Category{})

	if err != nil {
		return nil, err
	}

	storage := repository.NewCategory(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(module.PageSizeMax),
		),
	)
	service := usecase.NewCategory(storage, opts.EventBox, opts.ServiceHelper)
	serviceAPI := usecase_shared.NewCategoryAPI(storage, opts.ServiceHelper)
	*c = append(*c, http_v1.NewCategory(section, service, serviceImage, metaOrderBy))

	return serviceAPI, nil
}

func newUnitCategoryImage(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
) (*usecase.CategoryImage, error) {
	fileAPI, err := opts.S3Pool.Provider(opts.Cfg.ModulesSettings.CatalogCategory.Image.FileProvider)

	if err != nil {
		return nil, err
	}

	storage := repository.NewCategoryImage(opts.PostgresAdapter)
	service := usecase.NewCategoryImage(
		opts.Cfg.ModulesSettings.CatalogCategory.Image.BaseDir,
		storage,
		fileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)
	*c = append(*c, http_v1.NewCategoryImage(section, service))

	return service, nil
}
