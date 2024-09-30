package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/catalog/trademark/module"
	"github.com/mondegor/go-sample/internal/factory/catalog/trademark"
)

// CreateModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func CreateModule(ctx context.Context, opts trademark.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)

	if l, err := createUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.PrepareEachController(l, mrfactory.WithPermission(module.Permission))...)
	}

	return list, nil
}
