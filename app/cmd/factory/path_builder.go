package factory

import (
	"strings"

	"github.com/mondegor/go-sample/config"

	"github.com/mondegor/go-webcore/mrpath/placeholderpath"
)

// NewImageURLBuilder - создаёт объект placeholderpath.Builder.
func NewImageURLBuilder(cfg config.Config) (*placeholderpath.Builder, error) {
	return placeholderpath.New(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/")+
			"/"+
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.BaseURL, "/"),
		placeholderpath.Placeholder,
	)
}
