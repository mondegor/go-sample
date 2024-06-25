package filestation

import (
	"context"

	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/internal/factory/filestation"
)

// NewModuleOptions - comment func.
func NewModuleOptions(_ context.Context, opts app.Options) (filestation.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
	)
	if err != nil {
		return filestation.Options{}, err
	}

	return filestation.Options{
		UsecaseHelper:  opts.UsecaseErrorWrapper,
		RequestParser:  opts.RequestParsers.String,
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitImageProxy: filestation.UnitImageProxyOptions{
			FileAPI: fileAPI,
			BaseURL: opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL,
		},
	}, nil
}
