package factory

import (
	"go-sample/internal/modules"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
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
		RequestParsers: &factory.RequestParsers{
			String: opts.RequestParsers.String,
			Image: view_shared.NewParserImage(
				opts.RequestParsers.KeyInt32,
				opts.RequestParsers.String,
				opts.RequestParsers.Image,
			),
			Parser: view_shared.NewParser(
				opts.RequestParsers.Int64,
				opts.RequestParsers.ItemStatus,
				opts.RequestParsers.KeyInt32,
				opts.RequestParsers.SortPage,
				opts.RequestParsers.String,
				opts.RequestParsers.Validator,
			),
		},
		ResponseSender: opts.ResponseSender,

		CategoryAPI:  opts.CatalogCategoryAPI,
		OrdererAPI:   opts.OrdererAPI,
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
