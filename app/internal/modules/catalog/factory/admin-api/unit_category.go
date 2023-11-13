package factory

import (
	"go-sample/internal/factory"
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
	unitNameCategory = moduleName + ".Category"
)

func newUnitCategory(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
	storage *repository.Category,
	listSorter mrview.ListSorter,
) error {
	fileAPI, err := factory.NewS3MinioFileProvider(opts.MinioAdapter, opts.Cfg.BucketName, opts.Logger)

	if err != nil {
		return err
	}

	imageStorage := repository.NewCategoryImage(opts.PostgresAdapter)
	imageService := usecase.NewCategoryImage(
		opts.Cfg.FileStorage.CatalogCategoryImageDir,
		imageStorage,
		fileAPI,
		opts.Locker,
		opts.EventBox,
		opts.ServiceHelper,
	)

	service := usecase.NewCategory(storage, opts.EventBox, opts.ServiceHelper)
	*c = append(*c, http_v1.NewCategory(section, service, imageService, listSorter))

	return nil
}

func newUnitCategoryStorage(opts *modules.Options) (*repository.Category, *mrsql.EntityMetaOrderBy, error) {
	metaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity.Category{})

	if err != nil {
		return nil, nil, err
	}

	return repository.NewCategory(
		opts.PostgresAdapter,
		mrsql.NewBuilderSelect(
			mrpostgres.NewSqlBuilderWhere(),
			mrpostgres.NewSqlBuilderOrderByWithDefaultSort(metaOrderBy.DefaultSort()),
			mrpostgres.NewSqlBuilderPager(1000),
		),
	), metaOrderBy, nil
}
