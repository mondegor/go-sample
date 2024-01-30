package factory

import (
	"context"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(ctx context.Context, opts factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(ctx, module.Name)
	mrfactory.InfoCreateUnit(ctx, module.UnitCategoryName)

	if l, err := createUnitCategory(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitCategoryPermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitProductName)

	if l, err := createUnitProduct(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitProductPermission)...)
	}

	mrfactory.InfoCreateUnit(ctx, module.UnitTrademarkName)

	if l, err := createUnitTrademark(ctx, opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(ctx, l, module.UnitTrademarkPermission)...)
	}

	return list, nil
}
