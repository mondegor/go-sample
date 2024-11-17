package catalog

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/catalog/category/api/availability/usecase"
	"github.com/mondegor/go-sample/internal/catalog/category/shared/validate"
	"github.com/mondegor/go-sample/internal/factory/catalog/category"
	"github.com/mondegor/go-sample/internal/factory/catalog/category/api/availability"
)

// NewCategoryModuleOptions - создаёт объект category.Options.
func NewCategoryModuleOptions(_ context.Context, opts app.Options) (category.Options, error) {
	imageFileAPI, err := opts.FileProviderPool.ProviderAPI(
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
		EventEmitter:        opts.EventEmitter,
		UseCaseErrorWrapper: mrapp.NewUseCaseErrorWrapper(),
		DBConnManager:       opts.PostgresConnManager,
		Locker:              opts.Locker,
		RequestParsers: category.RequestParsers{
			// Parser:       opts.RequestParsers.Parser,
			// ExtendParser: opts.RequestParsers.ExtendParser,
			ModuleParser: validate.NewCategoryParser(
				opts.RequestParsers.ExtendParser,
				opts.RequestParsers.Image,
			),
		},
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

// NewCategoryAvailabilityAPI - создаёт объект usecase.Category.
func NewCategoryAvailabilityAPI(ctx context.Context, opts app.Options) (*usecase.Category, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init catalog category availability API")

	return availability.NewCategory(opts.PostgresConnManager), nil
}
