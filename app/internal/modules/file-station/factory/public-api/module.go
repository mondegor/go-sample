package factory

import (
	module "go-sample/internal/modules/file-station"
	"go-sample/internal/modules/file-station/factory"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
)

func CreateModule(opts *factory.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	mrfactory.InfoCreateModule(opts.Logger, module.Name)
	mrfactory.InfoCreateUnit(opts.Logger, module.UnitImageProxyName)

	if l, err := createUnitImageProxy(opts); err != nil {
		return nil, err
	} else {
		list = append(list, mrfactory.WithPermission(l, module.UnitImageProxyPermission)...)
	}

	return list, nil
}
