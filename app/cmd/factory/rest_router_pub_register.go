package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/go-sample/internal/app"
	catalogcategorypub "github.com/mondegor/go-sample/internal/factory/catalog/category/section/pub"
	filestationpub "github.com/mondegor/go-sample/internal/factory/filestation/section/pub"
)

// RegisterRestRouterPubHandlers - регистрирует в указанном роутере обработчики секции PublicAPI.
func RegisterRestRouterPubHandlers(ctx context.Context, router mrserver.HttpRouter, opts app.Options) error {
	section := NewAppSectionPublicAPI(ctx, opts)
	prepareHandler := mrfactory.WithMiddlewareCheckAccess(ctx, section, opts.AccessControl)
	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())

	for _, createFunc := range getPublicAPIControllers(ctx, opts) {
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

func getPublicAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return catalogcategorypub.CreateModule(ctx, opts.CatalogCategoryModule)
		},
		func() ([]mrserver.HttpController, error) {
			return filestationpub.CreateModule(ctx, opts.FileStationModule)
		},
	}
}
