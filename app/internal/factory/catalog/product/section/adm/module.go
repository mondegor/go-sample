package adm

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/product/module"
	"github.com/mondegor/go-sample/internal/factory/catalog/product"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

// CreateModule - comment func.
func CreateModule(ctx context.Context, opts product.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitProduct(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
