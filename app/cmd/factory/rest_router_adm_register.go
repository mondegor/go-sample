package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/go-sample/internal/app"
	catalogcategoryadm "github.com/mondegor/go-sample/internal/factory/catalog/category/section/adm"
	catalogproductadm "github.com/mondegor/go-sample/internal/factory/catalog/product/section/adm"
	catalogtrademarkadm "github.com/mondegor/go-sample/internal/factory/catalog/trademark/section/adm"
)

// RegisterRestRouterAdmHandlers - регистрирует в указанном роутере обработчики секции AdminAPI.
func RegisterRestRouterAdmHandlers(ctx context.Context, router mrserver.HttpRouter, opts app.Options) error {
	section := NewAppSectionAdminAPI(ctx, opts)
	prepareHandler := mrfactory.WithMiddlewareCheckAccess(ctx, section, opts.AccessControl)
	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())

	for _, createFunc := range getAdminAPIControllers(ctx, opts) {
		list, err := createFunc()
		if err != nil {
			return err
		}

		router.Register(
			mrfactory.PrepareEachController(list, prepareHandler)...,
		)
	}

	return nil
}

func getAdminAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return catalogcategoryadm.CreateModule(ctx, opts.CatalogCategoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogproductadm.CreateModule(ctx, opts.CatalogProductModule)
		},
		func() ([]mrserver.HttpController, error) {
			return catalogtrademarkadm.CreateModule(ctx, opts.CatalogTrademarkModule)
		},
	}
}
