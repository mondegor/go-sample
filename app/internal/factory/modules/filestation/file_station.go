package filestation

import (
	"context"

	"go-sample/internal/app"
	"go-sample/internal/modules/filestation/factory"
)

// NewModuleOptions - comment func.
func NewModuleOptions(_ context.Context, opts app.Options) (factory.Options, error) {
	fileAPI, err := opts.FileProviderPool.Provider(
		opts.Cfg.ModulesSettings.FileStation.ImageProxy.FileProvider,
	)
	if err != nil {
		return factory.Options{}, err
	}

	return factory.Options{
		UsecaseHelper:  opts.UsecaseErrorWrapper,
		RequestParser:  opts.RequestParsers.String,
		ResponseSender: opts.ResponseSenders.FileSender,

		UnitImageProxy: factory.UnitImageProxyOptions{
			FileAPI: fileAPI,
			BaseURL: opts.Cfg.ModulesSettings.FileStation.ImageProxy.BaseURL,
		},
	}, nil
}
