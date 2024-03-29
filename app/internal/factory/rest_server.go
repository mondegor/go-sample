package factory

import (
	"context"
	"go-sample/internal"
	factory_catalog_category_adm "go-sample/internal/modules/catalog/category/factory/admin-api"
	factory_catalog_category_pub "go-sample/internal/modules/catalog/category/factory/public-api"
	factory_catalog_product_adm "go-sample/internal/modules/catalog/product/factory/admin-api"
	factory_catalog_trademark_adm "go-sample/internal/modules/catalog/trademark/factory/admin-api"
	factory_filestation_pub "go-sample/internal/modules/file-station/factory/public-api"
	"time"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	restServerCaption = "RestServer"
)

func NewRestServer(ctx context.Context, opts app.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouter(ctx, opts.Cfg, opts.Translator)

	if err != nil {
		return nil, err
	}

	// section: admin-api
	sectionAdminAPI := NewAppSectionAdminAPI(ctx, opts)

	if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
		return nil, err
	}

	err = createAdminAPIControllers(
		ctx,
		opts,
		registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionAdminAPI, opts.AccessControl),
		),
	)

	if err != nil {
		return nil, err
	}

	// section: public
	sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

	if err = RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
		return nil, err
	}

	err = createPublicAPIControllers(
		ctx,
		opts,
		registerControllers(
			router,
			mrfactory.WithMiddlewareCheckAccess(ctx, sectionPublicAPI, opts.AccessControl),
		),
	)

	if err != nil {
		return nil, err
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrserver.NewServerAdapter(
		ctx,
		mrserver.ServerOptions{
			Caption:         restServerCaption,
			Handler:         router,
			ReadTimeout:     srvOpts.ReadTimeout * time.Second,
			WriteTimeout:    srvOpts.WriteTimeout * time.Second,
			ShutdownTimeout: srvOpts.ShutdownTimeout * time.Second,
			Listen: mrserver.ListenOptions{
				AppPath:  opts.Cfg.AppPath,
				Type:     srvOpts.Listen.Type,
				SockName: srvOpts.Listen.SockName,
				BindIP:   srvOpts.Listen.BindIP,
				Port:     srvOpts.Listen.Port,
			},
		},
	), nil
}

func registerControllers(router mrserver.HttpRouter, operations ...mrfactory.PrepareHandlerFunc) mrfactory.ApplyEachControllerFunc {
	return func(list []mrserver.HttpController, err error) error {
		if err != nil {
			return err
		}

		router.Register(
			mrfactory.PrepareEachController(list, operations...)...,
		)

		return nil
	}
}

func createAdminAPIControllers(ctx context.Context, opts app.Options, register mrfactory.ApplyEachControllerFunc) error {
	if err := register(factory_catalog_category_adm.CreateModule(ctx, opts.CatalogCategoryModule)); err != nil {
		return err
	}

	if err := register(factory_catalog_product_adm.CreateModule(ctx, opts.CatalogProductModule)); err != nil {
		return err
	}

	if err := register(factory_catalog_trademark_adm.CreateModule(ctx, opts.CatalogTrademarkModule)); err != nil {
		return err
	}

	return nil
}

func createPublicAPIControllers(ctx context.Context, opts app.Options, register mrfactory.ApplyEachControllerFunc) error {
	if err := register(factory_catalog_category_pub.CreateModule(ctx, opts.CatalogCategoryModule)); err != nil {
		return err
	}

	if err := register(factory_filestation_pub.CreateModule(ctx, opts.FileStationModule)); err != nil {
		return err
	}

	return nil
}
