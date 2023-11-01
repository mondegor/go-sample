package usecase

import (
    "context"
    "strings"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    FileItem struct {
        storage mrstorage.FileProvider
    }
)

func NewFileItem(
    storage mrstorage.FileProvider,
) *FileItem {
    return &FileItem{
        storage: storage,
    }
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *FileItem) Get(ctx context.Context, path string) (*mrentity.File, error) {
    path = strings.TrimLeft(path, "/")

    if path == "" {
        return nil, mrcore.FactoryErrServiceEmptyInputData.New("path")
    }

    item := mrentity.File{
        Name: path,
    }

    err := uc.storage.Download(ctx, &item)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, mrentity.ModelNameFile)
    }

    return &item, nil
}
