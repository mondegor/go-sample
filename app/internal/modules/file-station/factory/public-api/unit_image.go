package factory

import (
	"go-sample/internal/modules"
	http_v1 "go-sample/internal/modules/file-station/controller/http_v1/public-api"
	usecase "go-sample/internal/modules/file-station/usecase/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	unitNameImageProxy = moduleName + ".ImageProxy"
)

func newUnitImageProxy(
	c *[]mrcore.HttpController,
	opts *modules.Options,
	section mrcore.ClientSection,
) error {
	fileAPI, err := opts.S3Pool.Provider(opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider)

	if err != nil {
		return err
	}

	service := usecase.NewFileProviderAdapter(fileAPI, opts.ServiceHelper)
	*c = append(*c, http_v1.NewImageProxy(section, service, opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL))

	return nil
}
