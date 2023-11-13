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
    opts.Logger.Info("Init unit %s.%s in %s section", moduleName, unitNameCategory, section.Caption())

    return newUnitCategory(c, opts, section)
}
