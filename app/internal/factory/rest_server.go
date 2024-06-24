package factory

import (
	"context"

	"go-sample/internal/app"
	factory_catalog_category_adm "go-sample/internal/modules/catalog/category/factory/admin_api"
	factory_catalog_category_pub "go-sample/internal/modules/catalog/category/factory/public_api"
	factory_catalog_product_adm "go-sample/internal/modules/catalog/product/factory/admin_api"
	factory_catalog_trademark_adm "go-sample/internal/modules/catalog/trademark/factory/admin_api"
	factory_filestation_pub "go-sample/internal/modules/filestation/factory/public_api"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	restServerCaption = "RestServer"
)

// NewRestServer - comment func.
func NewRestServer(ctx context.Context, opts app.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouter(ctx, opts, opts.Translator)
	if err != nil {
		return nil, err
	}

	// section: admin-api
	{
		sectionAdminAPI := NewAppSectionAdminAPI(ctx, opts)

		if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
			return nil, err
		}

		registerControllersFunc := registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionAdminAPI, opts.AccessControl),
		)

		for _, createFunc := range getAdminAPIControllers(ctx, opts) {
			list, err := createFunc()
			if err != nil {
				return nil, err
			}

			registerControllersFunc(list)
		}
	}

	// section: public
	{
		sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

		if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
			return nil, err
		}

		registerControllersFunc := registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionPublicAPI, opts.AccessControl),
		)

		for _, createFunc := range getPublicAPIControllers(ctx, opts) {
			list, err := createFunc()
			if err != nil {
				return nil, err
			}

			registerControllersFunc(list)
		}
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrserver.NewServerAdapter(
		ctx,
		mrserver.ServerOptions{
			Caption:         restServerCaption,
			Handler:         router,
			ReadTimeout:     srvOpts.ReadTimeout,
			WriteTimeout:    srvOpts.WriteTimeout,
			ShutdownTimeout: srvOpts.ShutdownTimeout,
			Listen: mrserver.ListenOptions{
				BindIP: srvOpts.Listen.BindIP,
				Port:   srvOpts.Listen.Port,
			},
		},
	), nil
}

func registerControllers(router mrserver.HttpRouter, operations ...mrfactory.PrepareHandlerFunc) func(list []mrserver.HttpController) {
	return func(list []mrserver.HttpController) {
		router.Register(
			mrfactory.PrepareEachController(list, operations...)...,
		)
	}
}

func getAdminAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return factory_catalog_category_adm.CreateModule(ctx, opts.CatalogCategoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return factory_catalog_product_adm.CreateModule(ctx, opts.CatalogProductModule)
		},
		func() ([]mrserver.HttpController, error) {
			return factory_catalog_trademark_adm.CreateModule(ctx, opts.CatalogTrademarkModule)
		},
	}
}

func getPublicAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return factory_catalog_category_pub.CreateModule(ctx, opts.CatalogCategoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return factory_filestation_pub.CreateModule(ctx, opts.FileStationModule)
		},
	}
}
