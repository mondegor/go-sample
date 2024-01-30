package factory

import (
	"context"
	"go-sample/config"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlog"
)

func NewFileProviderPool(ctx context.Context, cfg config.Config) (*mrstorage.FileProviderPool, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init file provider pool")

	pool := mrstorage.NewFileProviderPool()

	//fs := NewFileSystem(ctx, cfg)
	//
	//if err := RegisterFileImageStorage(ctx, cfg, pool, fs); err != nil {
	//	return nil, err
	//}

	minioAdapter, err := NewS3Minio(ctx, cfg)

	if err != nil {
		return nil, err
	}

	if err = RegisterS3ImageStorage(ctx, cfg, pool, minioAdapter); err != nil {
		return nil, err
	}

	return pool, nil
}
