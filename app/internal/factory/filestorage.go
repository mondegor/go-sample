package factory

import (
    "errors"
    "go-sample/config"
    "os"

    "github.com/mondegor/go-storage/mrfilestorage"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewFileStorage(cfg *config.Config, logger mrcore.Logger) (mrstorage.FileProviderAPI, error) {
    logger.Info("Init file storage")

    err := createBaseDir(cfg.FileStorage.DownloadDir, 0755)

    if err != nil {
        return nil, err
    }

    err = os.MkdirAll(cfg.FileStorage.DownloadDir + "/" + cfg.FileStorage.CatalogCategoryImageDir, 0755)

    if err != nil {
        return nil, err
    }

    return mrfilestorage.New(cfg.FileStorage.DownloadDir), nil
}

func createBaseDir(dirPath string, perms os.FileMode) error {
    _, err := os.Stat(dirPath)

    if errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(dirPath, perms)
    }

    return err
}
