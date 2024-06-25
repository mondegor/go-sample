package pub

import (
	"context"

	"github.com/mondegor/go-sample/internal/factory/filestation"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/filestation/module"
)

// CreateModule - comment func.
func CreateModule(ctx context.Context, opts filestation.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitImageProxy(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.UnitImageProxyPermission))...)
	}

	return list, nil
}
