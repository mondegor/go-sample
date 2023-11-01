package factory

import (
    "go-sample/config"
    view_shared "go-sample/internal/controller/http_v1/shared/view"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrview"
)

func NewValidator(cfg *config.Config, logger mrcore.Logger) (mrcore.Validator, error) {
    logger.Info("Create and init data validator")

    validator := mrview.NewValidator()

    if err := validator.Register("article", view_shared.ValidateArticle); err != nil {
        return nil, err
    }

    return validator, nil
}
