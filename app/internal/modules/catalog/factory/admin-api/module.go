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
    categoryStorage, categoryMetaOrderBy, err := newUnitCategoryStorage(opts)

    if err != nil {
        return err
    }

    trademarkStorage, trademarkMetaOrderBy, err := newUnitTrademarkStorage(opts)

    if err != nil {
        return err
    }

    opts.Logger.Info("Init unit %s in %s section", unitNameCategory, section.Caption())

    if err = newUnitCategory(c, opts, section, categoryStorage, categoryMetaOrderBy); err != nil {
        return err
    }

    opts.Logger.Info("Init unit %s in %s section", unitNameTrademark, section.Caption())

    if err = newUnitTrademark(c, opts, section, trademarkStorage, trademarkMetaOrderBy); err != nil {
        return err
    }

    opts.Logger.Info("Init unit %s in %s section", unitNameProduct, section.Caption())

    if err = newUnitProduct(c, opts, section, categoryStorage, trademarkStorage); err != nil {
        return err
    }

    return nil
}
