package factory

import (
	"context"
	"go-sample/internal/modules"
	factory_catalog_adm "go-sample/internal/modules/catalog/factory/admin-api"
	factory_catalog_pub "go-sample/internal/modules/catalog/factory/public-api"
	factory_filestation_pub "go-sample/internal/modules/file-station/factory/public-api"
	"time"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

const (
	restServerCaption = "RestServer"
)

func NewRestServer(ctx context.Context, opts modules.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouter(ctx, opts.Cfg, opts.Translator)

	if err != nil {
		return nil, err
	}

	// section: admin-api
	sectionAdminAPI := NewAppSectionAdminAPI(ctx, opts)

	if err := RegisterSystemHandlers(ctx, opts.Cfg, router, sectionAdminAPI); err != nil {
		return nil, err
	}

	if controllers, err := factory_catalog_adm.CreateModule(ctx, opts.CatalogModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionAdminAPI, opts.AccessControl)...,
		)
	}

	// section: public
	sectionPublicAPI := NewAppSectionPublicAPI(ctx, opts)

	if err := RegisterSystemHandlers(ctx, opts.Cfg, router, sectionPublicAPI); err != nil {
		return nil, err
	}

	if controllers, err := factory_catalog_pub.CreateModule(ctx, opts.CatalogModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionPublicAPI, opts.AccessControl)...,
		)
	}

	if controllers, err := factory_filestation_pub.CreateModule(ctx, opts.FileStationModule); err != nil {
		return nil, err
	} else {
		router.Register(
			mrfactory.WithMiddlewareCheckAccess(ctx, controllers, sectionPublicAPI, opts.AccessControl)...,
		)
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
