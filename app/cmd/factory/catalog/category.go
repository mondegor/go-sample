package catalog

import (
	"context"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/category/api/availability/usecase"
	"github.com/mondegor/go-sample/internal/catalog/category/module"
	"github.com/mondegor/go-sample/internal/catalog/category/shared/validate"
	"github.com/mondegor/go-sample/internal/factory/catalog/category"
	"github.com/mondegor/go-sample/internal/factory/catalog/category/api/availability"
	"github.com/mondegor/go-sample/pkg/catalog/api"

	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewCategoryModuleOptions - comment func.
func NewCategoryModuleOptions(_ context.Context, opts app.Options) (category.Options, error) {
	imageFileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.CatalogCategory.Image.FileProvider,
	)
	if err != nil {
		return category.Options{}, err
	}

	categoryDictionary, err := opts.Translator.Dictionary("catalog/categories")
	if err != nil {
		return category.Options{}, err
	}

	return category.Options{
		EventEmitter:  opts.EventEmitter,
		UsecaseHelper: opts.UsecaseErrorWrapper,
		DBConnManager: opts.PostgresConnManager,
		Locker:        opts.Locker,
		RequestParser: validate.NewParser(
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

		UnitCategory: category.UnitCategoryOptions{
			Dictionary:      categoryDictionary,
			ImageFileAPI:    imageFileAPI,
			ImageURLBuilder: opts.ImageURLBuilder,
		},

		PageSizeMax:     opts.Cfg.General.PageSizeMax,
		PageSizeDefault: opts.Cfg.General.PageSizeDefault,
	}, nil
}

// NewCategoryAPI - comment func.
func NewCategoryAPI(ctx context.Context, opts app.Options) (*usecase.Category, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog category API")

	return availability.NewCategory(opts.PostgresConnManager, opts.UsecaseErrorWrapper), nil
}

// RegisterCategoryErrors - comment func.
func RegisterCategoryErrors(em *mrinit.ErrorManager) {
	em.RegisterList(mrinit.WrapProtoList(api.CategoryErrors()))
	em.RegisterList(mrinit.WrapProtoList(module.Errors()))
}
