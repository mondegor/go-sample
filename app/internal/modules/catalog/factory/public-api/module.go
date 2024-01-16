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
	opts.Logger.Info("Init unit %s.%s in %s section", module.Name, module.UnitCategoryName, section.Caption())

	return newUnitCategory(c, opts, section)
}
