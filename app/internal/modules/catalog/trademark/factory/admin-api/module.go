package factory

import (
	"context"
	module "go-sample/internal/modules/catalog/trademark"
	"go-sample/internal/modules/catalog/trademark/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.Permission)...)
	}

	return list, nil
}
