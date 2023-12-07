package factory

import (
	"go-sample/internal/modules"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	moduleName = "Catalog"
)

func NewModule(opts *modules.Options, section mrcore.ClientSection) ([]mrcore.HttpController, error) {
	opts.Logger.Info("Init module %s in section %s", moduleName, section.Caption())

	var c []mrcore.HttpController

	if err := newModule(&c, opts, section); err != nil {
		return nil, err
	}

	return c, nil
}

func newModule(c *[]mrcore.HttpController, opts *modules.Options, section mrcore.ClientSection) error {
	opts.Logger.Info("Init unit %s in %s section", unitNameCategory, section.Caption())

	categoryImage, err := newUnitCategoryImage(c, opts, section)

	if err != nil {
		return err
	}

	categoryAPI, err := newUnitCategory(c, opts, section, categoryImage)

	if err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", unitNameTrademark, section.Caption())

	trademarkAPI, err := newUnitTrademark(c, opts, section)

	if err != nil {
		return err
	}

	opts.Logger.Info("Init unit %s in %s section", unitNameProduct, section.Caption())

	if err = newUnitProduct(c, opts, section, categoryAPI, trademarkAPI); err != nil {
		return err
	}

	return nil
}
