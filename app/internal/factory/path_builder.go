package factory

import (
	"strings"

	"go-sample/config"

	"github.com/mondegor/go-webcore/mrpath/placeholderpath"
)

// NewImageURLBuilder - comment func.
func NewImageURLBuilder(cfg config.Config) (*placeholderpath.Builder, error) {
	return placeholderpath.New(
		strings.TrimRight(cfg.ModulesSettings.FileStation.ImageProxy.Host, "/")+
			"/"+
			strings.TrimLeft(cfg.ModulesSettings.FileStation.ImageProxy.BaseURL, "/"),
		placeholderpath.Placeholder,
	)
}
