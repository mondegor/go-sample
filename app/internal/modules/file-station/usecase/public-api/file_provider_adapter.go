package usecase

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAdapter struct {
		fileAPI       mrstorage.FileProviderAPI
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewFileProviderAdapter(
	fileAPI mrstorage.FileProviderAPI,
	serviceHelper *mrtool.ServiceHelper,
) *FileProviderAdapter {
	return &FileProviderAdapter{
		fileAPI:       fileAPI,
		serviceHelper: serviceHelper,
	}
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *FileProviderAdapter) Get(ctx context.Context, path string) (*mrtype.File, error) {
	path = strings.TrimLeft(path, "/")

	if path == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	file, err := uc.fileAPI.Download(ctx, path)

	if err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", path)
	}

	return file, nil
}
