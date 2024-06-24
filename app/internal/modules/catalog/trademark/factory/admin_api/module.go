package factory

import (
	"context"

	"go-sample/internal/modules/catalog/trademark/factory"
	"go-sample/internal/modules/catalog/trademark/module"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

// CreateModule - comment func.
func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
