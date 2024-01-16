package factory

import (
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/factory"

	"github.com/mondegor/go-webcore/mrcore"
)

func NewModule(opts *factory.Options, section mrcore.ClientSection) ([]mrcore.HttpController, error) {
	opts.Logger.Info("Init module %s in section %s", module.Name, section.Caption())

	var c []mrcore.HttpController

	if err := newModule(&c, opts, section); err != nil {
		return nil, err
	}

	return c, nil
}

func newModule(c *[]mrcore.HttpController, opts *factory.Options, section mrcore.ClientSection) error {
	opts.Logger.Info("Init unit %s in %s section", module.UnitCategoryName, section.Caption())

	if err := newUnitCategory(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitTrademarkName, section.Caption())

	if err := newUnitTrademark(c, opts, section); err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", module.UnitProductName, section.Caption())

	if err := newUnitProduct(c, opts, section); err != nil {
		return err
	}

	return nil
}
