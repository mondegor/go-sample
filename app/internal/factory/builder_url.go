package factory

import (
	"go-sample/config"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
)

func NewBuilderImagesURL(cfg *config.Config) mrcore.BuilderPath {
	return mrlib.NewBuilderPath(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/") +
			"/" +
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.URLPathRoot, "/") +
			"/",
	)
}
