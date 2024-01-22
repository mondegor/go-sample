package factory

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitCategoryName)

	if l, err := createUnitCategory(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitCategoryPermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitProductName)

	if l, err := createUnitProduct(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitProductPermission)...)
	}

	mrfactory.InfoCreateUnit(opts.Logger, module.UnitTrademarkName)

	if l, err := createUnitTrademark(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitTrademarkPermission)...)
	}

	return list, nil
}
