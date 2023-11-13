package usecase

import (
    "context"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtype"
)

type (
    FileProviderAdapter struct {
        fileAPI mrstorage.FileProviderAPI
    }
)

func NewFileProviderAdapter(
    fileAPI mrstorage.FileProviderAPI,
) *FileProviderAdapter {
    return &FileProviderAdapter{
        fileAPI: fileAPI,
    }
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *FileProviderAdapter) Get(ctx context.Context, path string) (*mrtype.File, error) {
    path = strings.TrimLeft(path, "/")

    if path == "" {
        return nil, mrcore.FactoryErrServiceEmptyInputData.New("path")
    }

    file, err := uc.fileAPI.Download(ctx, path)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, "FileProviderAPI")
    }

    return file, nil
}
