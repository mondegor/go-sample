package factory

import (
	"context"

	httpv1 "go-sample/internal/modules/filestation/controller/httpv1/public_api"
	"go-sample/internal/modules/filestation/factory"
	usecase "go-sample/internal/modules/filestation/usecase/public_api"

	"github.com/mondegor/go-webcore/mrserver"
)

func createUnitImageProxy(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitImageProxy(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

func newUnitImageProxy(_ context.Context, opts factory.Options) (*httpv1.ImageProxy, error) { //nolint:unparam
	useCase := usecase.NewFileProviderAdapter(opts.UnitImageProxy.FileAPI, opts.UsecaseHelper)
	controller := httpv1.NewImageProxy(
		opts.RequestParser,
		opts.ResponseSender,
		useCase,
		opts.UnitImageProxy.BaseURL,
	)

	return controller, nil
}
