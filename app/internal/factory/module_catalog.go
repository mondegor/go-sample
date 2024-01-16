package factory

import (
	"go-sample/internal/modules"
	"go-sample/internal/modules/catalog/factory"
	factory_api "go-sample/internal/modules/catalog/factory/api"
	usecase_api "go-sample/internal/modules/catalog/usecase/api"
)

func NewCatalogOptions(opts *modules.Options) (*factory.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.CatalogCategory.Image.FileProvider,
	)

	if err != nil {
		return nil, err
	}

	categoryDictionary, err := opts.Translator.Dictionary("catalog/category")

	if err != nil {
		return nil, err
	}

	return &factory.Options{
		Logger:          opts.Logger,
		EventBox:        opts.EventBox,
		ServiceHelper:   opts.ServiceHelper,
		PostgresAdapter: opts.PostgresAdapter,
		Locker:          opts.Locker,

		CategoryAPI:  opts.CatalogCategoryAPI,
		TrademarkAPI: opts.CatalogTrademarkAPI,

		UnitCategory: &factory.UnitCategoryOptions{
			Dictionary:      categoryDictionary,
			ImageFileAPI:    fileAPI,
			ImageURLBuilder: NewBuilderImagesURL(opts.Cfg),
		},
	}, nil
}

func NewCatalogCategoryAPI(opts *modules.Options) (*usecase_api.Category, error) {
	opts.Logger.Info("Create and init catalog category API")

	return factory_api.NewCategory(opts.PostgresAdapter, opts.ServiceHelper), nil
}

func NewCatalogTrademarkAPI(opts *modules.Options) (*usecase_api.Trademark, error) {
	opts.Logger.Info("Create and init catalog trademark API")

	return factory_api.NewTrademark(opts.PostgresAdapter, opts.ServiceHelper), nil
}
