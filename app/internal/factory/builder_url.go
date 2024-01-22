package factory

import (
	"go-sample/config"
	"strings"

	"github.com/mondegor/go-webcore/mrlib"
)

func NewBuilderImagesURL(cfg *config.Config) *mrlib.BuilderPath {
	return mrlib.NewBuilderPath(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/") +
			"/" +
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.BaseURL, "/") +
			"/",
	)
}
