package factory

import (
	"github.com/mondegor/go-webcore/mrcore/mrcoreerr"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-sample/config"
)

// NewErrorHandler - создаёт объект mrcoreerr.ErrorHandler.
func NewErrorHandler(_ mrlog.Logger, _ config.Config) *mrcoreerr.ErrorHandler {
	return mrcoreerr.NewErrorHandler()
}
