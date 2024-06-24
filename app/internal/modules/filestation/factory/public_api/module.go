package factory

import (
	"context"

	"go-sample/internal/modules/filestation/factory"
	"go-sample/internal/modules/filestation/module"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

// CreateModule - comment func.
func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitImageProxy(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.UnitImageProxyPermission))...)
	}

	return list, nil
}
