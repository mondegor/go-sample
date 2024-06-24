package catalog

import (
	"context"

	"go-sample/internal/modules/catalog/category/module"
	"go-sample/pkg/modules/catalog/api"

	"go-sample/internal/app"
	view_shared "go-sample/internal/modules/catalog/category/controller/httpv1/shared/view"
	"go-sample/internal/modules/catalog/category/factory"
	factory_api "go-sample/internal/modules/catalog/category/factory/api"
	usecase_api "go-sample/internal/modules/catalog/category/usecase/api"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewCategoryModuleOptions - comment func.
func NewCategoryModuleOptions(_ context.Context, opts app.Options) (factory.Options, error) {
	imageFileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.CatalogCategory.Image.FileProvider,
	)
	if err != nil {
		return factory.Options{}, err
	}

	categoryDictionary, err := opts.Translator.Dictionary("catalog/categories")
	if err != nil {
		return factory.Options{}, err
	}

	return factory.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		Locker:        opts.Locker,
		RequestParser: view_shared.NewParser(
			// opts.RequestParsers.Bool,
			// opts.RequestParsers.DateTime,
			opts.RequestParsers.Int64,
			// opts.RequestParsers.KeyInt32,
			opts.RequestParsers.ListSorter,
			opts.RequestParsers.ListPager,
			opts.RequestParsers.String,
			opts.RequestParsers.UUID,
			opts.RequestParsers.Validator,
			// opts.RequestParsers.File,
			opts.RequestParsers.Image,
			opts.RequestParsers.ItemStatus,
		),
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitCategory: factory.UnitCategoryOptions{
			Dictionary:      categoryDictionary,
			ImageFileAPI:    imageFileAPI,
			ImageURLBuilder: opts.ImageURLBuilder,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewCategoryAPI - comment func.
func NewCategoryAPI(ctx context.Context, opts app.Options) (*usecase_api.Category, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog category API")

	return factory_api.NewCategory(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterCategoryErrors - comment func.
func RegisterCategoryErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.CategoryErrors()))
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
