package factory

import (
	"go-sample/config"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewFileProviderPool(cfg *config.Config, logger mrcore.Logger) (*mrstorage.FileProviderPool, error) {
	logger.Info("Create and init file provider pool")

	pool := mrstorage.NewFileProviderPool()

	//fs := NewFileSystem(cfg, logger)
	//
	//if err := RegisterFileImageStorage(cfg, pool, fs, logger); err != nil {
	//	return nil, err
	//}

	minioAdapter, err := NewS3Minio(cfg, logger)

	if err != nil {
		return nil, err
	}

	if err = RegisterS3ImageStorage(cfg, pool, minioAdapter, logger); err != nil {
		return nil, err
	}

	return pool, nil
}
